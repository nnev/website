# encoding: utf-8



module Jekyll
	class LangPolyfill < Generator
		priority :high

		def generate(site)
			fakes = site.pages.map do |p|
				next if p.data['lang'] != 'de'
				next if File.exist?("en/#{p.name}")

				copy = p.dup
				copy.data['lang'] = 'en'
				copy.process("en/#{copy.name}")
				copy.data['title'] << " (untranslated, sorry)"

				notice = <<~NOTICE
				# Translation missing

				Sorry, unfortunately the current page has not been translated yet. [Maybe an
				automated translation suffices](https://translate.google.com/translate?hl=en&sl=de&tl=en&u=#{site.config['url']}#{copy.url})?
				NOTICE
				copy.content = "#{notice}\n\n#{copy.content}"

				copy
			end.compact

			site.pages += fakes
		end
	end
end
