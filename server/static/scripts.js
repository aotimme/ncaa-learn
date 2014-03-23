$("#predict-form").submit(function(event) {
  event.preventDefault();
  var home = $("#home-team").val();
  var away = $("#away-team").val();
  $.ajax({
    url: "/matchup",
    dataType: "json",
    data: { home: home, away: away},
    success: function(data, textStatus, jqhxr) {
      $("#winner").html(data.winner.name);
      $("#loser").html(data.loser.name);
      $("#score").html(parseInt(data.score.winner) + "-" + parseInt(data.score.loser));
      $("#probability").html(Math.round(data.score.probability * 1000)/10);
    }
  });
});

var teamSearcher = function(q, callback) {
  $.ajax({
    url: "/search",
    dataType: "json",
    data: { name: q },
    success: function(data) {
      var results = [];
      $.each(data, function(i, name) {
        results.push({value: name});
      });
      callback(results);
    }
  });
}

$('.typeahead.team-input').typeahead({
  hint: true,
  highlight: true,
  minLength: 1
}, {
  name: 'teams',
  displayKey: 'value',
  source: teamSearcher
});
