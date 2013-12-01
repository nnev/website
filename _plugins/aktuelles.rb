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
			conn = PGconn.open(:dbname => 'nnev')
			res = conn.exec('SELECT stammtisch, override, location, termine.date AS date, topic, abstract, vortraege.id AS c14h_id FROM termine LEFT JOIN vortraege ON termine.vortrag = vortraege.id ORDER BY termine.date LIMIT 4')
			termine = []
			res.each do |tuple|
				tuple['stammtisch'] = (tuple['stammtisch'] == 't')
				termine << tuple
			end

			res = conn.exec('SELECT * FROM vortraege ORDER BY date')
			vortraege = []
			res.each do |tuple|
				vortraege << tuple
			end

			site.pages.each do |page|
				page.data['termine'] = termine
				page.data['vortraege'] = vortraege
			end
		end
	end
end
