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

Meistens [Im Neuenheimer Feld 368, Raum 432](anfahrt.html). Jeden ersten Donnerstag
im Monat ist [Stammtisch](stammtisch.html) (abwechselnde Gaststätte).

Nächstes Treffen
===

{% assign termin = page.termine | first %}

<p itemscope itemtype="http://schema.org/Event">
	<time itemprop="startDate" datetime="{{termin.date}}T19:00"><b>{{ termin.date | escape }}</b> um 19 Uhr</time><br/>
	{% if termin.stammtisch %}
		<b itemprop="summary">Chaos-Stammtisch</b> bei
		{% for st in site.pages %}
			{% unless st.layout == "stammtisch" %}
				{% continue %}
			{% endunless %}
			{% unless st.name == termin.location %}
				{% continue %}
			{% endunless %}
			{% assign done = true %}
			<a href="{{ st.url | escape }}" itemprop="location">{{ termin.location | escape }}</a>
			<span itemprop="geo" itemscope itemtype="http://schema.org/GeoCoordinates">
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
    <span itemprop="geo" itemscope itemtype="http://schema.org/GeoCoordinates">
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
		<a href="pizza.html">Möchtest Du Pizza mitbestellen?</a>
	{% endif %}
</p>
