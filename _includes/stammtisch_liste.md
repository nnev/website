<ul>
{% for st in site.pages %}
  {% if st.layout == "stammtisch" %}
    <li><a href="{{st.url}}">
      {% if st.url == page.url %}
        <b>{{st.name}}</b>
      {% else %}
        {{st.name}}
      {% endif %}
    </a></li>
  {% endif %}
{% endfor %}
</ul>
