{% for termin in page.termine %}
<p {% if notfirst %}class="dim"{% endif %}>
	{% assign notfirst = true %}

	<b>Datum: {{ termin.date | escape }}</b><br>
	{% if termin.override != "" %}
		{{ termin.override | escape }}
	{% elsif termin.stammtisch %}
		<a href="stammtisch.html">Stammtisch</a><br/>
		<a href="yarpnarp.html">bitte zu/absagen</a><br/>
		{% if termin.location != "" %}
			{% for st in site.pages %}
				{% unless st.layout == "stammtisch" %}
					{% continue %}
				{% endunless %}
				{% unless st.name == termin.location %}
					{% continue %}
				{% endunless %}
				{% assign done = true %}
				Location: <a href="{{ st.url }}">{{ termin.location | escape }}</a>
			{% endfor %}
			{% unless done %}
				Location: <a href="stammtisch.html">{{ termin.location | escape }}</a>
			{% endunless %}
		{% else %}
			Location: TBA
		{% endif %}
	{% else %}
		Chaos-Treff (<a href="anfahrt.html">Anfahrt?</a>)<br>
		c¼h:
		{% if termin.topic %}
			<a href="chaotische_viertelstunde.html#c14h_{{termin.c14h_id}}">{{ termin.topic | escape }}</a>
		{% else %}
			noch keine ◉︵◉.<br/>
			<a href="edit_c14.html?date={{termin.date | escape}}">neue c¼h eintragen?</a>
		{% endif %}
	{% endif %}
</p>
{% endfor %}
