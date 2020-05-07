# encoding: utf-8

module Jekyll
  module LinkifyFilter
    def linkify(input)
      re = Regexp.new "(https://[^ ]+)", Regexp::MULTILINE
      input.gsub re, '<a href="\1" data-no-turbolink>\1</a>'
    end
  end
end

Liquid::Template.register_filter(Jekyll::LinkifyFilter)
