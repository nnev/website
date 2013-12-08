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
im Monat ist Stammtisch (abwechselnde Gaststätte).

Nächstes Treffen
===

{% assign termin = page.termine | first %}

<p>
  <b>{{ termin.date | escape }}</b> um 19 Uhr<br/>
  {% if termin.stammtisch %}
    <b>Stammtisch</b> bei
    {% for st in site.pages %}
      {% unless st.layout == "stammtisch" %}
        {% continue %}
      {% endunless %}
      {% unless st.name == termin.location %}
        {% continue %}
      {% endunless %}
      {% assign done = true %}
      <a href="{{ st.url | escape }}">{{ termin.location | escape }}</a>
    {% endfor %}
    {% unless done %}
      <a href="stammtisch.html">{{ termin.location | escape }}</a>
    {% endunless %}
    <br>
    <a href="FIXME">Zwecks Reservierung bitte zu/absagen</a>
  {% else %}
    <b>Chaos-Treff</b> <a href="anfahrt.html">(Anfahrt)</a><br/>
    c¼h:
    {% if termin.topic %}
      <a href="chaotische_viertelstunde.html#c14h_{{termin.c14h_id}}">{{ termin.topic | escape }}</a>
    {% else %}
      noch keine ◉︵◉
    {% endif %}
    <br>
    <a href="pizza.html">Möchtest Du Pizza mitbestellen?</a>
  {% endif %}
</p>
