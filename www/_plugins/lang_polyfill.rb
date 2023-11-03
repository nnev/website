# encoding: utf-8



module Jekyll
	class LangPolyfill < Generator
		priority :high

		def generate(site)
			fakes = site.pages.map do |p|
				next if p.data['lang'] != 'de'
				next if File.exist?("en/#{p.path}")

				copy = p.dup
				copy.data = p.data.dup

				copy.data['lang'] = 'en'
				copy.process("en/#{copy.path}")

				copy.data['title'] << " (untranslated, sorry)"
				copy.data['translation_missing'] = <<~NOTICE
				<p class="translation_missing">
					<b>Translation missing</b>.	Unfortunately the current page has not been translated yet. Maybe an <a href="https://translate.google.com/translate?hl=en&sl=de&tl=en&u=#{site.config['url']}#{copy.url}">automated translation</a> suffices?
				</p>
				NOTICE

				copy
			end.compact

			site.pages += fakes
		end
	end
end
