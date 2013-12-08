---
layout: default
title: Zu-/Absagen für den Chaos-Stammtisch
---

# Chaos-Stammtisch

Erscheinst Du zum nächsten Stammtisch? Gib hier bitte Deine ja/nein
Stimme ab, damit ich passend reservieren kann. Danke!

### Zu-/Absagen für Chaos-Stammtisch

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
	<input type="text" placeholder="Dein Nick" id="nick" name="nick" value="<<.Nick>>" /><br>

	<label for="kommentar">Kommentar</label>
	<input type="text" placeholder="Kommentar" id="kommentar" name="" value="<<.Kommentar>>" /><br>

	<input type="submit" value="Yarp" name="kommt"/>
	<input type="submit" value="Narp" name="kommt"/>
</form>


### Status
