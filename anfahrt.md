---
layout: default
title: Anfahrt zum Chaos-Treff Heidelberg
---

# Anfahrt


<a href="https://maps.google.de/maps?q=49.417433,8.675255&num=1&t=m&z=18" class="qrcode">
	<img src="/img/map_google.png"/><br/>
	Google Maps
</a>

<a href="http://www.openstreetmap.org/?mlat=49.41737&mlon=8.67527#map=17/49.41737/8.67527" class="qrcode"  style="margin-left: 2rem">
	<img src="/img/map_osm.png"/><br/>
	OSM.org
</a>

<address>
Erdgeschoss, Seminarraum A
Mathematikon
Im Neuenheimer Feld 205
69120 Heidelberg
</address>

Am besten kommt man zu uns, wenn man von Norden in den Innenhof läuft und an der Tür in der Süd-West-Ecke kurz klopft.

Die nächste Haltestelle ist das „[Bun&shy;sen&shy;gym&shy;na&shy;sium](http://fahrplanauskunft.vrn.de/vrn/XSLT_TRIP_REQUEST2?language=de&sessionID=0&name_destination=Neuenheim,%20Bunsengymnasium&type_destination=stop)“ (Bus 31, 32, Tram 21, 24). Par&shy;ken ist im Stadtteil Neuenheim (auf der anderen Seite der Berliner Straße) kostenfrei möglich, vor 19 Uhr sollte aber eine Parkscheibe verwendet werden [\[1\]](https://www.heidelberg.de/site/Heidelberg_ROOT/get/documents/heidelberg/PB5Documents/pdf/81_pdf_Parken_Neuenheim_1uebersichtsplan.pdf). Alternativ gibt es im nördlichen Bauteil des Mathematikons eine kostenpflichtige Tiefgarage.

<div id="map"></div>
<script>
var map = L.map('map').setView([49.417433, 8.675255], 16);
{{site.map_js}}
L.marker([{{site.treff_lat}}, {{site.treff_lon}}]).addTo(map).bindPopup("<b>Chaos-Treff</b>", { "closeButton": false }).openPopup();
L.marker([49.41694, 8.67633]).addTo(map).bindPopup("Bus- und Tramhaltestelle: Bunsengymnasium");
</script>
