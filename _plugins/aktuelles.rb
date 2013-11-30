# encoding: utf-8

module Jekyll
	class Aktuelles < Jekyll::Generator
		def generate(site)
			site.pages.each do |page|

				page.data['termine'] = [
					{
						"date" => "2013-12-26",
						"stammtisch" => false,
						"topic" => "FÄLLT AUS (30c3)"
					}, {
						"date" => "2013-12-19",
						"stammtisch" => false,
						"topic" => "Penisverlängerungen"
					}, {
						"date" => "2013-12-12",
						"stammtisch" => false,
						"topic" => "This stock is about to go through the roof!"
					}, {
						"date" => "2012-12-05",
						"stammtisch" => true,
						"location" => "Mister Wu"
					}
				]
			end
		end
	end
end
