$(function() {
  $("#predict-form").submit(function(event) {
    event.preventDefault();
    var home = $("#home-team").val();
    var away = $("#away-team").val();
    $.ajax({
      url: "/matchup",
      dataType: "json",
      data: {home: home, away: away},
      success: function(data, textStatus, jqhxr) {
        $("#winner").html(data.winner.name);
        $("#loser").html(data.loser.name);
        $("#score").html(parseInt(data.score.winner) + "-" + parseInt(data.score.loser));
        $("#probability").html(Math.round(data.score.probability * 1000)/10);
      }
    });
  });

  var teamSearcher = function(q, syncResults, asyncResults) {
    $.ajax({
      url: "/search",
      dataType: "json",
      data: {name: q},
      success: function(data) {
        var results = $.map(data, function(name) {
          return {value: name};
        });
        asyncResults(results);
      },
      error: function(jqxhr) {
        console.error('ERROR', jqxhr);
      }
    });
  };

  $('.typeahead.team-input').typeahead({
    hint: false,
    highlight: true,
    minLength: 1
  }, {
    name: 'teams',
    async: true,
    display: 'value',
    limit: 10,
    source: teamSearcher
  });
});