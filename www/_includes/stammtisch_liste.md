<ul>
{% for st in site.pages %}
  {% if st.layout == "stammtisch" %}
    <li><a href="{{st.url | escape}}">
      {% if st.url == page.url %}
        <b>{{st.locname | escape}}</b>
      {% else %}
        {{st.locname | escape}}
      {% endif %}
    </a></li>
  {% endif %}
{% endfor %}
</ul>
