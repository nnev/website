---
layout: default
title: NoName e.V.
---

Was?
===

Treff zum Austausch über Computer, Technik und deren gesellschaftliche Auswirkungen

Wann?
===

Jeden Donnerstag ab 19 Uhr.

Wo?
===

Meistens [Im Neuenheimer Feld 368, Raum 432](anfahrt.html).<br/>
Jeden ersten Donnerstag im Monat ist [Stammtisch](stammtisch.html) (abwechselnde Gaststätte).

rgb2r?
===

Am kommenden Wochenende (2014-10-31 - 2014-11-02) findet die „rgb2r – roots go back to the roots“, unser Retro-Gaming-Event statt. Die Anmeldung ist zwar schon geschlossen aber du kannst einfach vorbei kommen. [Alle Einzelheiten dazu findest du auf der rgb2r Seite](http://rgb2r.noname-ev.de/).

Nächstes Treffen
===

{% assign termin = page.termine | first %}

<p itemscope itemtype="http://data-vocabulary.org/Event">
	<time itemprop="startDate" datetime="{{termin.date}}T19:00"><b>{{ termin.date | escape }}</b> um 19 Uhr</time><br/>
	{% if termin.override != "" %}
		{{ termin.override | escape }}
	{% elsif termin.stammtisch %}
		<b itemprop="summary">Chaos-Stammtisch</b> bei
		{% for st in site.pages %}
			{% unless st.layout == "stammtisch" %}
				{% continue %}
			{% endunless %}
			{% unless st.title == termin.location %}
				{% continue %}
			{% endunless %}
			{% assign done = true %}
			<a href="{{ st.url | escape }}" itemprop="location">{{ termin.location | escape }}</a>
			<span itemprop="geo" itemscope itemtype="http://data-vocabulary.org/Geo">
				<meta itemprop="latitude" content="{{st.lat}}" />
				<meta itemprop="longitude" content="{{st.lon}}" />
			</span>
		{% endfor %}
		{% unless done %}
			<a href="stammtisch.html" itemprop="location">{{ termin.location | escape }}</a>
		{% endunless %}
		<br>
		<a href="yarpnarp.html" itemprop="url">Zwecks Reservierung bitte zu/absagen</a>
	{% else %}
		<b itemprop="summary">Chaos-Treff</b> <a href="anfahrt.html">(Anfahrt)</a><br/>
    <span itemprop="geo" itemscope itemtype="http://data-vocabulary.org/Geo">
      <meta itemprop="latitude" content="{{site.treff_lat}}" />
      <meta itemprop="longitude" content="{{site.treff_lon}}" />
    </span>
		c¼h:
		{% if termin.topic %}
			<a itemprop="url" href="chaotische_viertelstunde.html#c14h_{{termin.c14h_id}}">{{ termin.topic | escape }}</a>
		{% else %}
			noch keine ◉︵◉
		{% endif %}
		<br>
		<a href="pizza.html">Möchtest Du Pizza mitbestellen?</a><br>Hinweis: Diesen Donnerstag (30.10.) ist die Pizza-Bestellung <b>nur bis 17 Uhr</b> offen und wird gegen 19 Uhr erwartet, damit die anschließende Mitgliederversammlung reibungsfrei ablaufen kann.
	{% endif %}
</p>
