<ul>
{% for st in site.pages %}
  {% if st.layout == "stammtisch" %}
    <li><a href="{{st.url | escape}}">
      {% if st.url == page.url %}
        <b>{{st.title | escape}}</b>
      {% else %}
        {{st.title | escape}}
      {% endif %}
    </a></li>
  {% endif %}
{% endfor %}
</ul>
