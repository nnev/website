---
layout: default
title: Zu-/Absagen für den Chaos-Stammtisch
---

# Chaos-Stammtisch

Erscheinst Du zum nächsten Stammtisch? Gib hier bitte Deine ja/nein
Stimme ab, damit wir passend reservieren können. Danke!

<strong>Aufgepasst:</strong> Manchmal müssen wir den Stammtisch spontan in ein anderes
Lokal verlegen, weil im geplanten nicht genug Platz für uns ist. Bitte schau daher 
nochmal hier auf diese Seite, bevor du zu Hause losläufst!

### Zu-/Absagen

<p>
{% for termin in page.termine %}
	{% if termin.stammtisch %}
		<b>Datum: {{ termin.date }}</b><br>
		{% if termin.location != "" %}
			{% for st in site.pages %}
				{% unless st.layout == "stammtisch" %}
					{% continue %}
				{% endunless %}
				{% unless st.locname == termin.location %}
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
		{% break %}
	{% endif %}
{% endfor %}
</p>


<form method="POST">
	<label for="nick">Dein Nick</label>
	<input type="text" placeholder="Dein Nick" id="nick" name="nick" value="<<.Nick>>" required="required"/><br>

	<label for="kommentar">Kommentar</label>
	<input type="text" placeholder="Kommentar (Optional)" id="kommentar" name="kommentar" value="<<.Kommentar>>" /><br>

	<label for="captcha">Captcha: Gebe hier "NoName e.V." ein:</label>
	<input type="text" placeholder="NoName e.V." id="captcha" name="captcha" required="required" /><br>

	<input type="submit" value="Yarp" name="kommt"/>
	<input type="submit" value="Narp" name="kommt"/>
</form>


### Status

<div>
<< if .HasKommt >>
	<< if .Kommt >>
		<div class="yarpnarpstatus yarp">
			<b>Dein Status:</b> Du kommst (ツ)
		</div>
	<< else >>
		<div class="yarpnarpstatus narp">
			<b>Dein Status:</b> Du kommst nicht (⊙︿⊙)
		</div>
	<< end >>
<< end >>
</div>

<table class="yarpnarp">
	{% if page.stammtischYarpCount != 0 %}
		<tr class="header"><th>Yarp</th><td>{{page.stammtischYarpCount}} Zusagen</td></tr>
		{% for yarp in page.stammtischYarp %}
			<tr>
				<td>{{yarp.nick | escape}}</td>
				<td title="{{yarp.kommentar | escape}}">{{yarp.kommentar | escape}}</td>
			</tr>
		{% endfor %}
	{% endif %}

	{% if page.stammtischNarpCount != 0 %}
		<tr class="header"><th>Narp</th><td>{{page.stammtischNarpCount}} Absagen</td></tr>
		{% for narp in page.stammtischNarp %}
			<tr>
				<td>{{narp.nick | escape}}</td>
				<td title="{{narp.kommentar | escape}}">{{narp.kommentar | escape}}</td>
			</tr>
		{% endfor %}
	{% endif %}
</table>
