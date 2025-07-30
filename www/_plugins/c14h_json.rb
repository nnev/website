# encoding: utf-8

require 'json'
require 'pg'

module Jekyll
	class C14hJSON < Generator
		def generate(site)
			return real(site) if ENV['DONT_HIDE_FAILURES']

			begin
				real(site)
			rescue => e
				warn "\n\nc14h-JSON Plugin ist kaputt. Keine Datenbank vorhanden? Fehlermeldung:"
				warn e.message
				warn e.backtrace.map{|x| "\t#{x}"}.join("\n")
				warn "\n\n"
			end
		end

		def real(site)
			conn = PG.connect(:dbname => 'nnev')
			vortraege = conn.exec('SELECT id,date,topic,speaker,abstract FROM vortraege ORDER BY date DESC').to_a
			vortraege.each do | vortrag |
				linklist = conn.exec('SELECT kind, url FROM vortrag_links WHERE vortrag = %s' % vortrag["id"]).to_a
				vortrag["links"] = linklist if linklist.length > 0
				vortrag["id"] = vortrag["id"].to_i
				vortrag.delete("date") if not vortrag["date"]
				vortrag.delete("abstract") if vortrag["abstract"].length == 0
			end

			FileUtils.mkdir_p(site.dest)
			File.open(File.join(site.dest, "c14h.json"), "w") do |f|
				f.write(JSON.generate(vortraege))
			end

			site.static_files << Jekyll::SitemapFile.new(site, site.dest, "/", "c14h.json")
		end
	end
end
