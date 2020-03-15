# encoding: utf-8

require "pg"

module Jekyll
	class Aktuelles < Jekyll::Generator
		def generate(site)
			return real(site) if ENV['DONT_HIDE_FAILURES']

			begin
				real(site)
			rescue => e
				warn "\n\nAktuelles ist kaputt. Fehlermeldung:"
				warn e.message
				warn e.backtrace.map{|x| "\t#{x}"}.join("\n")
				warn "\n\n"
			end
		end

		def real(site)
			conn = PG.connect(:dbname => 'nnev')
			res = conn.exec('SELECT stammtisch, override, override_long, location, termine.date AS date, topic, abstract, vortraege.id AS c14h_id FROM termine LEFT JOIN vortraege ON termine.date = vortraege.date WHERE termine.date >= CURRENT_DATE ORDER BY termine.date LIMIT 4')
			termine = []
			res.each do |tuple|
				tuple['stammtisch'] = (tuple['stammtisch'] == 't')
				termine << tuple
			end

			vortraege_zukunft = conn.exec('SELECT * FROM vortraege WHERE date >= CURRENT_DATE ORDER BY date ASC').to_a

			vortraege_tba = conn.exec('SELECT * FROM vortraege WHERE date IS NULL ORDER BY id ASC').to_a

			vortraege_vergangenheit = conn.exec('SELECT * FROM vortraege WHERE date < CURRENT_DATE ORDER BY date DESC').to_a

			vortraege_vergangenheit.each do | vortrag |
				linklist = conn.exec('SELECT kind, url FROM vortrag_links WHERE vortrag = %s' % vortrag["id"]).to_a
				vortrag["links"] = linklist
			end

			vortraege_zukunft.each do | vortrag |
				linklist = conn.exec('SELECT kind, url FROM vortrag_links WHERE vortrag = %s' % vortrag["id"]).to_a
				vortrag["links"] = linklist
			end

			vortraege_tba.each do | vortrag |
				linklist = conn.exec('SELECT kind, url FROM vortrag_links WHERE vortrag = %s' % vortrag["id"]).to_a
				vortrag["links"] = linklist
			end


			latest = conn.exec('SELECT id FROM vortraege ORDER BY id DESC LIMIT 1')

			site.pages.each do |page|
				page.data['termine'] = termine
				page.data['vortraege_zukunft'] = vortraege_zukunft
				page.data['vortraege_tba'] = vortraege_tba
				page.data['vortraege_vergangenheit'] = vortraege_vergangenheit
				page.data['vortraege_latest'] = latest.getvalue(0, 0)
			end
		end
	end
end
