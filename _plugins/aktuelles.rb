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

			pp termine

			site.pages.each do |page|
				page.data['termine'] = termine
			end
		end
	end
end
