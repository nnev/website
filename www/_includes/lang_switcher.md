<div class="language_switcher">
  {% if page.lang == "de" %}
    <a href="{{ 'en/' | append: page.url | replace: '//','/' }}" lang="en">view page in English</a>
  {% else %}
    <a href="{{ page.url | replace: 'en/','' }}" lang="de">Seite auf Deutsch zeigen</a>
  {% endif %}

  {{ page.translation_missing }}
</div>

