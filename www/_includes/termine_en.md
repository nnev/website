{% for termin in page.termine %}
<p {% if notfirst %}class="dim"{% endif %}>
	{% assign notfirst = true %}

	<b>date (yyyy-mm-dd): {{ termin.date | escape }}</b> from 19h<br>
	{% if termin.override != "" %}
		{{ termin.override | escape }}
	{% elsif termin.stammtisch %}
		<a href="stammtisch.html">Restaurant (“Stammtisch”)</a><br/>
		<a href="yarpnarp.html">please RSVP</a><br/>
		{% if termin.location != "" %}
			{% for st in site.pages %}
				{% unless st.layout == "stammtisch" %}
					{% continue %}
				{% endunless %}
				{% unless st.locname == termin.location %}
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
		regular meeting (<a href="anfahrt.html">Route?</a>)<br>
		(short) talk:
		{% if termin.topic %}
			<a href="chaotische_viertelstunde.html#c14h_{{termin.c14h_id}}">{{ termin.topic | escape }}</a>
		{% else %}
			none yet ◉︵◉.
		{% endif %}
	{% endif %}
</p>
{% endfor %}

