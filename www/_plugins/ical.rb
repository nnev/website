# encoding: utf-8

require 'pg'
require 'pp'
require 'icalendar'
require 'date'
require 'fileutils'

include Icalendar

# enable unicode for icalendar.
#$KCODE = 'u'

module Jekyll
	class C14h < Generator
		priority :low

		def generate(site)
			return real(site) if ENV['DONT_HIDE_FAILURES']

			begin
				real(site)
			rescue => e
				warn "\n\nical-Plugin ist kaputt. Keine Datenbank vorhanden? Fehlermeldung:"
				warn e.message
				warn e.backtrace.map{|x| "\t#{x}"}.join("\n")
				warn "\n\n"
			end
		end

		def real(site)
			cal = Calendar.new
			cal.timezone do |t|
				t.tzid = 'Europe/Berlin'
			end

			conn = PGconn.open(:dbname => 'nnev')
			res = conn.exec(
				<<-SQL
				SELECT stammtisch, override, override_long, location, termine.date AS date, topic, abstract, speaker, vortraege.id AS c14h_id
				FROM termine
				LEFT JOIN vortraege
				ON termine.date = vortraege.date
				WHERE termine.date > now() - '1 years'::interval
				ORDER BY termine.date ASC
				SQL
			)

			stammtischs = site.pages.reject { |p| p.data['layout'] != "stammtisch" }

			res.each do |tuple|
				stammtisch = tuple['stammtisch'] == 't'
				desc = ""
				status = "TENTATIVE"
				if tuple['override'] != ""
					topic    = "NoName e.V.: #{tuple['override']}"
					desc     = tuple['override_long'] || ""
					status   = "CANCELLED"
				elsif stammtisch
					url      = "#{site.config['url']}/yarpnarp.html"

					details = stammtischs.find { |s| s.data['locname'] == tuple["location"] }

					topic    = "Chaos-Stammtisch"
					topic   << ": #{tuple["location"]}" unless tuple["location"].empty?

					desc << "=====\nbitte zu/absagen: #{url}\n=====\n\n"
					if details
						desc    << "#{details.content}\n\n"
						desc    << "#{details.data['menu_url']}\n#{details.data['site_url']}"
						desc    << "\n\n#{details.data['phone']}"
						desc    << "\n\n#{details.data['address']}\n#{details.data['address_desc']}"
						desc    << "\nhttp://www.openstreetmap.org/?mlat=#{details.data['lat']}&mlon=#{details.data['lon']}"
						desc    << "\n#{details.data['gmaps_url']}"
						status   = "CONFIRMED"
					end
					location  = tuple['location'].empty? ? "TBA" : tuple['location']
					location << ", #{details.data['address']}" if details && !details.data['address'].empty?
				else
					topic    = "Chaos-Treff"
					topic   << ": #{tuple['topic']}" if tuple['topic']
					status   = "CONFIRMED" if tuple['topic']

					desc     = "#{tuple['topic']}\n\n#{tuple['abstract']}"
					desc    << "\n\nVortragende/r: #{tuple['speaker']}" unless tuple['speaker'].nil? || tuple['speaker'].empty?
					desc    << "\n\nhttp://www.openstreetmap.org/?mlat=#{site.config['treff_lat']}&mlon=#{site.config['treff_lon']}"

					location = 'Im Neuenheimer Feld 368, Heidelberg'
					url      = "#{site.config['url']}/anfahrt.html"
				end

				cal.event do |e|
					e.dtstart     = DateTime.parse(tuple['date'] + ' 19:00')
					e.dtend       = DateTime.parse(tuple['date'] + ' 23:00')
					e.summary     = topic
					e.description = desc.strip
					e.organizer   = 'treff@noname-ev.de'
					e.location    = location
					e.status      = status
					e.uid         = "chaos-#{tuple['date']}@noname-ev.de"
					e.url         = url
				end
			end

			cal.publish

			FileUtils.mkdir_p(site.dest)
			File.open(File.join(site.dest, "c14h.ics"), "w") do |f|
				f.write(cal.to_ical)
			end

			site.static_files << Jekyll::SitemapFile.new(site, site.dest, "/", "c14h.ics")
		end
	end
end
