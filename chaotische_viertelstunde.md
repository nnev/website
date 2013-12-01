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
		<th><a id="c14h_{{vortrag.id}}" class="anchorhack"></a>{{ vortrag.date }}</th>
		<th colspan="2">{{ vortrag.topic }}</th>
	</tr>
	<tr><td></td><td class="dim" colspan="2">von {{ vortrag.speaker }}</td></tr>
	<tr>
		<td></td>
		<td class="just" colspan="2">{{ vortrag.abstract | newline_to_br }}</td>
	</tr>
	<tr class="space">
		<td></td>
		<td>
			{% if vortrag.infourl != '' %}
				<a href="{{ vortrag.infourl }}"><b>Details / Folien aufrufen</b></a>
			{% endif %}
		</td>
		<td>
			<a href="edit_c14.html?id={{ vortrag.id }}">Eintrag bearbeiten</a>
		</td>
	</tr>

{% endfor %}
</table>
