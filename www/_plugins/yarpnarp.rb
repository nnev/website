# encoding: utf-8

require "pg"

module Jekyll
	class YarpNarp < Jekyll::Generator
		def generate(site)
			return real(site) if ENV['DONT_HIDE_FAILURES']

			begin
				real(site)
			rescue => e
				warn "\n\nYarpNarp ist kaputt. Fehlermeldung:"
				warn e.message
				warn e.backtrace.map{|x| "\t#{x}"}.join("\n")
				warn "\n\n"
			end
		end

		def real(site)
			conn = PGconn.open(:dbname => 'nnev')

			yarp = conn.exec('SELECT nick, kommentar FROM zusagen WHERE kommt = TRUE ORDER BY nick ASC').to_a
      narp = conn.exec('SELECT nick, kommentar FROM zusagen WHERE kommt = FALSE ORDER BY nick ASC').to_a

			site.pages.each do |page|
				page.data['stammtischYarp'] = yarp
				page.data['stammtischYarpCount'] = yarp.size
				page.data['stammtischNarp'] = narp
				page.data['stammtischNarpCount'] = narp.size
			end
		end
	end
end
