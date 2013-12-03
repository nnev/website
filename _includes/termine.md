{% for termin in page.termine %}
<p {% if notfirst %}class="dim"{% endif %}>
	{% assign notfirst = true %}

	<b>Datum: {{ termin.date }}</b><br>
	{% if termin.stammtisch %}
		Stammtisch<br/>
		{% for st in site.pages %}
			{% unless st.layout == "stammtisch" %}
				{% continue %}
			{% endunless %}
			{% unless st.name == termin.location %}
				{% continue %}
			{% endunless %}
			{% assign done = true %}
			Location: <a href="{{ st.url }}">{{ termin.location }}</a>
		{% endfor %}
		{% unless done %}
			Location: <a href="stammtisch.html">{{ termin.location }}</a>
		{% endunless %}
	{% else %}
		Chaos-Treff (<a href="anfahrt.html">Anfahrt?</a>)<br>
		c¼h:
		{% if termin.topic %}
			<a href="chaotische_viertelstunde.html#c14h_{{termin.c14h_id}}">{{ termin.topic }}</a>
		{% else %}
			noch keine ◉︵◉.<br/>
			<a href="edit_c14.html">neue c¼h eintragen?</a>
		{% endif %}
	{% endif %}
</p>
{% endfor %}
