---
layout: default
title: Zu-/Absagen für den Chaos-Stammtisch
---

# Chaos-Stammtisch

Erscheinst Du zum nächsten Stammtisch? Gib hier bitte Deine ja/nein
Stimme ab, damit ich passend reservieren kann. Danke!

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
				{% unless st.name == termin.location %}
					{% continue %}
				{% endunless %}
				{% assign done = true %}
				Location: <a href="{{ st.url }}">{{ termin.location }}</a>
			{% endfor %}
			{% unless done %}
				Location: <a href="stammtisch.html">{{ termin.location }}</a>
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

	<input type="submit" value="Yarp" name="kommt"/>
	<input type="submit" value="Narp" name="kommt"/>
</form>


### Status

{% if HasKommt %}
	{% if Kommt %}
		<div class="yarpnarpstatus yarp">
			<b>Dein Status:</b> Du kommst (ツ)
		</div>
	{% else %}
		<div class="yarpnarpstatus narp">
			<b>Dein Status:</b> Du kommst nicht (⊙︿⊙)
		</div>
	{% endif %}
{% endif %}

<table class="yarpnarp">
	{% if page.stammtischYarpCount != 0 %}
		<tr class="header"><th>Yarp</th><td>{{page.stammtischYarpCount}} Zusagen</td></tr>
		{% for yarp in page.stammtischYarp %}
			<tr>
				<td>{{yarp.nick | escape}}</td>
				<td title="{{yarp.kommentar}}">{{yarp.kommentar}}</td>
			</tr>
		{% endfor %}
	{% endif %}

	{% if page.stammtischNarpCount != 0 %}
		<tr class="header"><th>Narp</th><td>{{page.stammtischNarpCount}} Absagen</td></tr>
		{% for narp in page.stammtischNarp %}
			<tr>
				<td>{{narp.nick | escape}}</td>
				<td title="{{narp.kommentar}}">{{narp.kommentar}}</td>
			</tr>
		{% endfor %}
	{% endif %}
</table>
