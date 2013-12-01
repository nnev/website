---
layout: default
title: Chaotische Viertelstunde
---

Seit knapp 2 Jahren findet bei jedem Treff ein Kurzvortrag statt, die so
gennante Chaotische Viertelstunde. Dabei gibt es wenige Regeln, aber
idealerweise geht der Vortrag ca. 15 Minuten. Ob du frei vorträgst, dein
Notebook oder die Tafel benutzt, bleibt dir überlassen.

Bei der Themenwahl hast du freie Hand. Themen, die wir in der Vergangenheit
hatten, drehten sich oft um Programmiersprachen (z.B. Go, CHICKEN Scheme),
Programme (sup, notmuch, Ingress, tor, sieve), Hardware (Raspberry Pi, Mifare
Classic, Human Enhancements), Life Hacking und vieles mehr.

Grundsätzlich gilt: uns gefällt alles! Du musst nicht nachfragen, ob es genug
Interessenten gibt — <a href="edit_c14.html">trag deinen Vortrag einfach ein</a>.

# Die nächsten Vorträge

<table>
{% for vortrag in page.vortraege %}
	<tr>
		<td>{{ vortrag.date }}</td>
		<td>{% if vortrag.infourl != '' %}<a href="{{ vortrag.infourl }}">{{ vortrag.topic }}</a>{% else %}{{ vortrag.topic }}{% endif %}</td>
		<td>{{ vortrag.speaker }}</td>
		<td><a href="edit_c14.html?id={{ vortrag.id }}">edit</a></td>
	</tr>
	<tr class="space">
		<td colspan="3">{{ vortrag.abstract }}</td>
	</tr>

{% endfor %}
</table>
