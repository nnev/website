---
layout: default
title: Chaos-Stammtisch
---

# Chaos-Stammtisch

Der Chaos-Stammtisch findet in wechselnden Lokalen statt. Wir reservieren nach
Bedarf, bitte nutze daher das [Anmeldesystem](yarpnarp.html). Der Treffpunkt
steht üblicherweise eine Woche im Voraus fest, ab dann kannst Du Dich auch ins
System eintragen. Wir starten so gegen 19 Uhr. Die Reservierung läuft im
Allgemeinen auf den Namen "F. Nord".

Du willst nichts verpassen? [Abonniere den ICS-Kalender](c14h.ics).

<div itemscope itemtype="http://data-vocabulary.org/Event">
<h3>Nächster <span itemprop="summary">Chaos-Stammtisch</span></h3>

{% for termin in page.termine %}
	{% assign notfirst = true %}
	{% if termin.override != "" %}
	{% elsif termin.stammtisch %}
		<p {% if notfirst %}class="dim"{% endif %}>
		<b>Datum: <time itemprop="startDate" datetime="{{termin.date}}T19:00">{{ termin.date | escape }}</time></b><br>
		<a href="yarpnarp.html" itemprop="url">bitte zu/absagen</a><br/>
		{% if termin.location != "" %}
			{% for st in site.pages %}
				{% unless st.layout == "stammtisch" %}
					{% continue %}
				{% endunless %}
				{% unless st.locname == termin.location %}
					{% continue %}
				{% endunless %}
				{% assign done = true %}
				Location: <a href="{{ st.url }}" itemprop="location">{{ termin.location | escape }}</a>
        <span itemprop="geo" itemscope itemtype="http://data-vocabulary.org/Geo">
          <meta itemprop="latitude" content="{{st.lat}}" />
          <meta itemprop="longitude" content="{{st.lon}}" />
        </span>
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
