<script type="text/javascript">
// for all
$(document).ready(function() {
  $('#value,#values').val("");
  $(".navbar-burger").click(function() {
    $(".navbar-burger").toggleClass("is-active");
    $(".navbar-menu").toggleClass("is-active");
  });
});

// for memo
$('#submit').click(function() {
  $('#submit').addClass('is-loading');
  var textType = "#value"
  if($('#toggle').is(':checked')) {
    textType = "#values"
  }
  var value = $(textType).val();
  if (value.length >= 200) {
    $('#warning').css("color", "#f14668").text("Must be 200 characters or less.")
    $('#submit').removeClass('is-loading');
  } else {
    $.ajax({
      type: "POST",
      url: "/insert",
      data: {
        "msg": value
      },
      success: function(result){
        location.reload();
      },
      error: function(xhr, status, error) {
        $('#submit').removeClass('is-loading');
        $('#warning').text($.parseJSON(xhr.responseText).msg);
      }
    });
  }
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

// for file
$('#upload').click(function() {
  $('#upload').addClass('is-loading');
  let file = $('#upload_file').prop('files')[0];
  // 3MiB
  let max = 3 * 1024 * 1024
  if (file.size > max) {
    $('#file_name').css("color", "#f14668").text("File size must be 3MiB or less.")
    $('#upload').removeClass('is-loading');
  } else {
    let reader = new FileReader();
    reader.addEventListener("load", function (){
      let fd = new FormData();
      fd.append("base64str", reader.result);
      fd.append("filename", file.name);
      $.ajax({
        url:'/upload',
        type:'post',
        data: fd,
        processData: false,
        contentType: false,
        cache: false,
        success: function(result) {
          location.reload();
        },
        error: function(xhr, status, error) {
          $('#upload').removeClass('is-loading');
          $('#msg').text($.parseJSON(xhr.responseText).msg);
        }
      });
    }, false);
    if (file) {
      reader.readAsDataURL(file);
    }
  }
});
$('#upload_file').change(function() {
    $('#file_name').css("color", "#ffffff")
  $('#upload').removeAttr('disabled');
  let file = $('#upload_file').prop('files')[0];
  $('#file_name').text(file.name);
});
$('#clear_file').click(function() {
  $('#clear_file').addClass('is-loading');
  $.ajax({
    type: "POST",
    url: "/clear_file",
    success: function(result) {
      location.reload();
    }
  });
});

// for user
$('#logout').click(function() {
  $.ajax({
    type: "POST",
    url: "/logout",
    success: function(result) {
      window.location.href = '/';
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
$('#save').click(function() {
  $('#save').addClass('is-loading');
  // Error parsing request body:invalid semicolon separator in query
  var img = board.getImg().replace('data:image/png;base64,', '');
  $.ajax({
    type: "POST",
    url: "/save",
    enctype: 'application/x-www-form-urlencoded',
    data: "image=" + img,
    success: function(msg){
      $('#save').removeClass('is-loading');
      window.location.href = '/image';
    },
    error: function(xhr, status, error) {
      $('#save').removeClass('is-loading');
      $('#msg').text($.parseJSON(xhr.responseText).msg);
    }
  });
});
$('#clear_img').click(function() {
  $('#clear_img').addClass('is-loading');
  $.ajax({
    type: "POST",
    url: "/clear_img",
    success: function(result) {
      $('#clear_img').removeClass('is-loading');
      window.location.href = '/image';
    }
  });
});
</script>
