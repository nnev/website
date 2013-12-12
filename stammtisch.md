---
layout: default
title: Chaos-Stammtisch
---

# Chaos-Stammtisch

Der Chaos-Stammtisch findet in wechselnden Lokalen statt. Wir reservieren nach Bedarf, bitte nutze daher das [Anmeldesystem FIXME LINK](FIXME). Der Treffpunkt steht üblicherweise eine Woche im Voraus fest, ab dann kannst Du Dich auch ins System eintragen.

### Nächster Stammtisch

<div>
{% for termin in page.termine %}
	{% assign notfirst = true %}
	{% if termin.override != "" %}
	{% elsif termin.stammtisch %}
		<p {% if notfirst %}class="dim"{% endif %}>
		<b>Datum: {{ termin.date | escape }}</b><br>
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
		</p>
		{% break %}
	{% endif %}
{% endfor %}
</div>

### Bisherige Locations

{% include stammtisch_liste.md %}
