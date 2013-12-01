# encoding: utf-8

require 'pg'
require 'pp'
require 'icalendar'
require 'date'
require 'fileutils'

include Icalendar

def parse_into_utc(datetime)
  DateTime.parse(datetime).new_offset(Rational(0, 24))
end

module Jekyll
  class C14h < Jekyll::Generator
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
      # Regardless of the time zone in which the host machine is running,
      # Chaostreffs always take place in Europe/Berlin time, so temporarily
      # switch to that to get the correct offset.
      prev_tv = ENV['TZ']
      ENV['TZ'] = 'Europe/Berlin'
      offset = DateTime.now.strftime('%z')
      ENV['TZ'] = prev_tv

      cal = Calendar.new
      cal.timezone do
        timezone_id 'UTC'
      end

			conn = PGconn.open(:dbname => 'nnev')
			res = conn.exec('SELECT stammtisch, override, location, termine.date AS date, topic, abstract FROM termine LEFT JOIN vortraege ON termine.vortrag = vortraege.id ORDER BY termine.date LIMIT 4')
			termine = []
			res.each do |tuple|
        location = if tuple['stammtisch'] == 't'
          # XXX: We could parse the lat/long from the stammtisch page.
          'Stammtisch'
        else
          'Im Neuenheimer Feld 368, Heidelberg'
        end

        cal.event do
          dtstart   parse_into_utc(tuple['date'] + ' 19:00'+offset)
          dtend     parse_into_utc(tuple['date'] + ' 23:59'+offset)
          summary   tuple['topic']
          organizer 'ccchd@ccchd.de'
          location  location
        end
			end

      cal.publish

      FileUtils.mkdir(site.dest)
      File.open(File.join(site.dest, "c14h.ics"), "w") do |f|
        f.write(cal.to_ical)
      end

      site.static_files << Jekyll::SitemapFile.new(site, site.dest, "/", "c14h.ics")
    end
  end
end
