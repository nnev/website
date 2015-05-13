---
layout: default
title: Chaotische Viertelstunde
foto: c14h
---

Seit Anfang 2012 findet bei jedem Treff ein Kurzvortrag statt, die so
genannte Chaotische Viertelstunde. Dabei gibt es wenige Regeln, aber
idealerweise geht der Vortrag ca. 15 Minuten. Ob du frei vorträgst, dein
Notebook oder die Tafel benutzt, bleibt dir überlassen.

Bei der Themenwahl hast du freie Hand. Themen, die wir in der Vergangenheit
hatten, drehten sich oft um Programmiersprachen (z.B. Go, CHICKEN Scheme),
Programme (sup, notmuch, Ingress, tor, sieve), Hardware (Raspberry Pi, Mifare
Classic, Human Enhancements), Life Hacking und vieles mehr.

Grundsätzlich gilt: uns gefällt alles! Du musst nicht nachfragen, ob es genug
Interessenten gibt — [trag deinen Vortrag einfach ein](edit_c14.html).

Du willst nichts verpassen? [Abonniere den ICS-Kalender](c14h.ics).

# Die nächsten Vorträge

{% assign vortraege = page.vortraege_zukunft %}
{% include c14h.html %}


<a class="button" onclick="document.getElementById('vortraege_tba').style.display='block';this.style.display='none';">Vorschläge ohne Datum anzeigen</a>
<div id="vortraege_tba">
{% assign vortraege = page.vortraege_tba %}
{% include c14h.html %}
</div>

# Das hast du verpasst

{% assign vortraege = page.vortraege_vergangenheit %}
{% include c14h.html %}
