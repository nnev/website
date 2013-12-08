
{% if page.foto %}
  <br/>
  <br/>
  <a href="/img/foto_{{page.foto | escape}}.jpg">
    <img src="/img/foto_{{page.foto | escape}}_thumb.jpg"/>
  </a>
{% endif %}
