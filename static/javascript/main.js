var expanded = false

function showMessage(level, message) {
  $("#messages").empty()
  var html =
  `<div class="alert alert-${level}" role="alert">
    <strong>${level}! </strong>${message}
  </div>`
  $("#messages").append(html)
}

function showMessages(messages) {
  $("#messages").empty()
  console.log(messages)
  $.each(messages, function(index, value) {
    var html =
    `<div class="alert alert-${value.severity}" role="alert">
      <strong>${value.severity}! </strong>${value.content}
    </div>`
    $("#messages").append(html)
  })
}

function create() {
  var action = $("#action option:selected").text()
  showSpinner()
  $.ajax({
    url: `/api/${action}`,
    type: "post",
    data : $("form").serialize(),
    success: function(data) {
      console.log(JSON.stringify(data))
      showMessages(data)
    },
    error: function(req, status, err) {
      if (req.responseJSON != undefined) {
        showMessage("danger", req.responseJSON.error)
      }
      else {
        showMessage("danger", JSON.stringify(req))
      }
    },
    complete: function() {
      hideSpinner()
    }
  })
}

function showList() {
  showSpinner()
  $.ajax({
    url: "/api/list/" + $('select#database option:selected').text(),
    type: "get",
    dataType: "json",
    success: function(data) {
      var list = $("#database-list")
      $.each(data, function(index, element) {
        list.append(`<small class="text-muted">${element.username}<br></small>`)
      })
    },
    complete: function() {
      hideSpinner()
      $("#card").show("slow")
      expanded = true
    }
  })
}

function hideList() {
  $( "#card" ).hide(function() {
    $(".card small").remove()
    expanded = false
  })
}

function toogleList() {
  expanded ? hideList() : showList()
}

function showSpinner() {
  $("#circle").removeClass('d-none')
}

function hideSpinner() {
  $("#circle").addClass('d-none')
}
