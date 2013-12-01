# encoding: utf-8

require "pg"
require "pp"

module Jekyll
	class Aktuelles < Jekyll::Generator
		def generate(site)
			conn = PGconn.open(:dbname => 'nnev')
			res = conn.exec('SELECT stammtisch, override, location, termine.date AS date, topic, abstract FROM termine LEFT JOIN vortraege ON termine.vortrag = vortraege.id ORDER BY termine.date LIMIT 4')
			termine = []
			res.each do |tuple|
				tuple['stammtisch'] = (tuple['stammtisch'] == 't')
				termine << tuple
			end

			res = conn.exec('SELECT * FROM vortraege')
			vortraege = []
			res.each do |tuple|
				vortraege << tuple
			end

			pp termine
			pp vortraege

			site.pages.each do |page|
				page.data['termine'] = termine
				page.data['vortraege'] = vortraege
			end
		end
	end
end
