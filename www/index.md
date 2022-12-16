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

Meistens [Mathematikon, Im Neuenheimer Feld 205, Seminarraum A](anfahrt.html).<br/>
Jeden ersten Donnerstag im Monat ist [Stammtisch](stammtisch.html) (abwechselnde Gaststätte).

Neugierig?
===

Komm einfach vorbei! Der Treff ist offen für alle, neue Gesichter sind immer gerne gesehen.

COVID-19 Schutzmaßnahmen
===

Da wir uns in den Räumlichkeiten der Universität treffen, gelten die
[Maßnahmen der Universität zum Schutz vor dem Coronavirus](https://www.uni-heidelberg.de/de/newsroom/massnahmen-der-universitaet-zum-schutz-vor-dem-coronavirus).


Nächstes Treffen
===

{% assign termin = page.termine | first %}

<p itemscope itemtype="http://data-vocabulary.org/Event">
	<time itemprop="startDate" datetime="{{termin.date}}T19:00"><b>{{ termin.date | escape }}</b> um 19 Uhr</time><br/>
	{% if termin.override != "" %}
		{{ termin.override | escape | linkify }}<br>
		{{ termin.override_long | markdownify }}
	{% elsif termin.stammtisch %}
		<b itemprop="summary">Chaos-Stammtisch</b> bei
		{% for st in site.pages %}
			{% unless st.layout == "stammtisch" %}
				{% continue %}
			{% endunless %}
			{% unless st.locname == termin.location %}
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
		<a href="pizza.html">Möchtest Du Pizza mitbestellen?</a>
	{% endif %}
</p>
