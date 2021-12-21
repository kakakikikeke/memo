<script type="text/javascript">
// for memo
$(function() {
  $(document).ready(function() {
    $('#value,#values').val("");
    $(".navbar-burger").click(function() {
      $(".navbar-burger").toggleClass("is-active");
      $(".navbar-menu").toggleClass("is-active");
    });
  });

  $('#submit').click(function() {
    $('#submit').addClass('is-loading');
    var textType = "#value"
    if($('#toggle').is(':checked')) {
      textType = "#values"
    }
    var value = $(textType).val();
    $.ajax({
      type: "POST",
      url: "/insert",
      data: {
        "msg": value
      },
      success: function(result){
        location.reload();
      }
    });
  });
  $('#clear').click(function() {
    $('#clear').addClass('is-loading');
    $.ajax({
      type: "POST",
      url: "/clear",
      success: function(result) {
        location.reload();
      }
    });
  });
  $('#check').click(function() {
    $('#check').addClass('is-loading');
    var name = $('#name').val();
    var pass = $('#pass').val();
    $.ajax({
      type: "POST",
      url: "/check",
      data: {
        "name": name,
        "pass": pass,
      },
      success: function(result) {
        $('#check').removeClass('is-loading');
        window.location.href = '/';
      },
      error: function(xhr, status, error) {
        $('#check').removeClass('is-loading');
        $('#msg').text($.parseJSON(xhr.responseText).msg);
      }
    });
  });
  $('#logout').click(function() {
    $.ajax({
      type: "POST",
      url: "/logout",
      success: function(result) {
        window.location.href = '/';
      }
    });
  });
  $('#create').click(function() {
    $('#create').addClass('is-loading');
    var name = $('#name').val();
    var pass = $('#pass').val();
    var pass2 = $('#pass2').val();
    $.ajax({
      type: "POST",
      url: "/create",
      data: {
        "name": name,
        "pass": pass,
        "pass2": pass2,
      },
      success: function(result) {
        $('#create').removeClass('is-loading');
        window.location.href = '/';
      },
      error: function(xhr, status, error) {
        $('#create').removeClass('is-loading');
        $('#msg').text($.parseJSON(xhr.responseText).msg);
      }
    });
  });
  $('#delete').click(function() {
    $('#delete').addClass('is-loading');
    $.ajax({
      type: "POST",
      url: "/delete",
      success: function(result) {
        $('#delete').removeClass('is-loading');
        window.location.href = '/';
      },
      error: function(xhr, status, error) {
        $('#delete').removeClass('is-loading');
        window.location.href = '/';
      }
    });
  });

  $('#value,#values').keyup(function() {
    if($(this).val().length != 0) {
      $('#submit').attr('disabled', false);
    } else {
      $('#submit').attr('disabled', true);
    }
  });
  $('#value').keypress(function(e) {
    if (e.which == 13)  {
      $('#submit').click();
    }
  });
  $('#toggle').click(function() {
    if($('#toggle').is(':checked')) {
      $('#value').css('display', 'none');
      $('#values').removeAttr('style');
    } else {
      $('#value').removeAttr('style');
      $('#values').css('display', 'none');
    }
  });
});
// for board
var board = new DrawingBoard.Board('board', {
  background: "#ffffff",
  color: "#000000",
  size: 30,
  fillTolerance: 150,
  controls: [
    { Size: { type: "range", min: 12, max: 42 } },
    { Navigation: { back: true, forward: true } },
    'DrawingMode',
    'Color'
  ],
  webStorage: 'local',
  droppable: true
});
board.addControl('Download');
board.downloadImg = function() {
  var img = this.getImg();
  img = img.replace("image/png", "image/octet-stream");
  var link = document.createElement('a');
  link.download = "download.png";
  link.href = img;
  link.click();
};
</script>