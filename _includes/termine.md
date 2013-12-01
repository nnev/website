{% for termin in page.termine %}
<p>
	Datum: {{ termin.date }}<br>
	{% if termin.stammtisch %}
	Stammtisch</br>
	{% for st in site.pages %}
		{% unless st.layout == "stammtisch" %}
			{% continue %}
		{% endunless %}
		{% unless st.name == termin.location %}
			{% continue %}
		{% endunless %}
		{% assign done = true %}
		Location: <a href="{{ st.link }}">{{ termin.location }}</a>
	{% endfor %}
	{% unless done %}
		Location: {{ termin.location }}
	{% endunless %}
	{% else %}
	<a href="anfahrt.html">Treff</a><br>
	c¼h: {{ termin.topic }}
	{% endif %}
</p>
{% endfor %}
