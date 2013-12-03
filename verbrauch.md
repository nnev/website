---
layout: default
title: Getränke-Verbrauch
---

# Getränke-Verbrauch

Wir erfassen unseren Getränke-Verbrauch digital. So richtig wissen wir noch nicht, was das bringt… aber in der Zwischenzeit haben wir einen tollen Graphen:

<script type="text/javascript" src="https://www.google.com/jsapi"></script>
<script type="text/javascript" src="https://noname-drink.appspot.com/stats"></script>
<script type="text/javascript">
  // Load the Visualization API and the piechart package.
  google.load('visualization', '1.0', {'packages':['corechart']});

  // Set a callback to run when the Google Visualization API is loaded.
  google.setOnLoadCallback(drawChart);

  // Callback that creates and populates a data table,
  // instantiates the pie chart, passes in the data and
  // draws it.
  function drawChart() {
var data = google.visualization.arrayToDataTable(nonamedrinkstats);
    // Set chart options
    var options = {'title':'Getränkeverbrauch',
                   'width': "100%",
                   'height':400,
       'curveType':'function'};

    // Instantiate and draw our chart, passing in some options.
    var chart = new google.visualization.LineChart(document.getElementById('chart_div'));
    chart.draw(data, options);
  }
</script>

<div id="chart_div"></div>
