
function updateKnownWifiList() {
  $.ajax({                                                                   
    type: "GET",                                                                        
    url: "/wifi",  
    contentType: "application/json; charset=utf-8",                                                            
    dataType: "json",   
    success: function(data) {
      var html ='';
      console.log("Got known wifilist:" + JSON.stringify(data));
      $.each(data, function(index, item) {
        html += '<li><a href="#"><h3>' + item.ssid+ '</h3><p>Priority:'+ item.priority + '</p></a><a href="#" data-icon="delete">Delete</a></li>';
      });
      $('#ul_knownwifinetworks').html($(html));
      $('#ul_knownwifinetworks').trigger('create');    
      $('#ul_knownwifinetworks').listview('refresh');
      
    },                                               
    error: function(msg) {              
      alert(msg.statusText);
    } 
  });  
  
}


function updateWifiList() {
  $.ajax({                                                                   
    type: "GET",                                                                        
    url: "/wifilist",  
    contentType: "application/json; charset=utf-8",                                                            
    dataType: "json",   
    success: function(data) {
      var html ='';
      console.log("Got wifilist:" + JSON.stringify(data));
      $.each(data, function(index, item) {
        html += '<li data-icon="plus"><a href="#"><h3>' + item.ssid+ '</h3><p>BSSID:'+ item.bssid + '</p></a></li>';
      });
      $('#ul_wifinetworks').html($(html));
      $('#ul_wifinetworks').trigger('create');    
      $('#ul_wifinetworks').listview('refresh');
      
    },                                               
    error: function(msg) {              
      alert(msg.statusText);
    } 
  });  
  
}




$(document).on("pageinit", "#configuration-page", function() {

  $(document).on('click',"#configuration-gear2", function () {
    $(".configuration-tab").hide();
    $("#wifi-tab").show();
    $("#wifi-navbar").addClass("ui-btn-active");
    $("#nabto-navbar").removeClass("ui-btn-active");
    $("#zs-navbar").removeClass("ui-btn-active");
  });


  updateWifiList();
  
  $(".configuration-tab").hide();
  $("#wifi-tab").show();
  updateKnownWifiList();
  $("#addwifibutton").on('click',function () {
    $(".configuration-tab").hide();
    $("#wifiaddnetwork-tab").show();
  });

  $("#nabto-navbar").on('click', function () {
    $(".configuration-tab").hide();
    $("#nabto-tab").show();
  });
  $("#wifi-navbar").on('click', function () {
    $(".configuration-tab").hide();
    updateKnownWifiList();
    
    $("#wifi-tab").show();
  });
  $("#zs-navbar").on('click', function () {
    $(".configuration-tab").hide();
    $("#zs-tab").show();
  });

  
});

