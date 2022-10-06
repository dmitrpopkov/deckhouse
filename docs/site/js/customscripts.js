$('#mysidebar').height($(".nav").height());

$( document ).ready(function() {
    $('#search-input').on("keyup", function (e) {
            if (e.target.value.length > 0 ) $(".search__results").addClass("active");
            else $(".search__results").removeClass("active");
          });
});

document.addEventListener("DOMContentLoaded", function() {
    /**
    * AnchorJS
    */
    if (window.anchors_disabled != true) {
        anchors.add('h2,h3,h4,h5');
        anchors.add('.anchored');
    }

});

$( document ).ready(function() {
    //this script says, if the height of the viewport is greater than 800px, then insert affix class, which makes the nav bar float in a fixed
    // position as your scroll. if you have a lot of nav items, this height may not work for you.
    var h = $(window).height();
    //console.log (h);
    if (h > 800) {
        // $( "#mysidebar" ).attr("class", "nav affix");
        // $( "#mysidebar" ).attr("class", "nav");
    }
    // activate tooltips. although this is a bootstrap js function, it must be activated this way in your theme.
    $('[data-toggle="tooltip"]').tooltip({
        placement : 'top'
    });

});

// needed for nav tabs on pages. See Formatting > Nav tabs for more details.
// script from http://stackoverflow.com/questions/10523433/how-do-i-keep-the-current-tab-active-with-twitter-bootstrap-after-a-page-reload
$(function() {
    var json, tabsState;
    $('a[data-toggle="pill"], a[data-toggle="tab"]').on('shown.bs.tab', function(e) {
        var href, json, parentId, tabsState;

        tabsState = localStorage.getItem("tabs-state");
        json = JSON.parse(tabsState || "{}");
        parentId = $(e.target).parents("ul.nav.nav-pills, ul.nav.nav-tabs").attr("id");
        href = $(e.target).attr('href');
        json[parentId] = href;

        return localStorage.setItem("tabs-state", JSON.stringify(json));
    });

    tabsState = localStorage.getItem("tabs-state");
    json = JSON.parse(tabsState || "{}");

    $.each(json, function(containerId, href) {
        return $("#" + containerId + " a[href=" + href + "]").tab('show');
    });

    $("ul.nav.nav-pills, ul.nav.nav-tabs").each(function() {
        var $this = $(this);
        if (!json[$this.attr("id")]) {
            return $this.find("a[data-toggle=tab]:first, a[data-toggle=pill]:first").tab("show");
        }
    });
});

$(document).ready(function() {
    var $notice = $('#notice');
    var $notice_collapse = $('#notice-collapse');
    var $notice_expand = $('#notice-expand');
    var notice_state = localStorage.getItem('notice-state') || 'expanded';

    function switchNotice(state) {
        $notice.attr('data-state', state);
        localStorage.setItem('notice-state', state);
    }

    switchNotice(notice_state);

    $notice_collapse.on('click', (e) => {
        e.preventDefault();
        switchNotice('collapsed');
    })
    $notice_expand.on('click', (e) => {
        e.preventDefault();
        switchNotice('expanded')
    })
});

/* features tabs */

$(document).ready(function() {
    $('[data-features-tabs-trigger]').on('click', function() {
        var name = $(this).attr('data-features-tabs-trigger');
        var $parent = $(this).closest('[data-features-tabs]');
        var $triggers = $parent.find('[data-features-tabs-trigger]');
        var $contents = $parent.find('[data-features-tabs-content]');
        var $content = $parent.find('[data-features-tabs-content=' + name + ']');

        $triggers.removeClass('active');
        $contents.removeClass('active');

        $(this).addClass('active');
        $content.addClass('active');
    })
});

/* Request access */
var ra = {};
ra.api_url = 'https://license.deckhouse.io/api/license/request';

function raSend(e) {
    e.preventDefault();
    if ($('#h0n3y').val() != '') {
        raError()
    } else {
        $.ajax({
            type: 'POST',
            url: ra.api_url,
            data: ra.form.serialize(),
            dataType: 'json',
            success: raSuccess,
            error: raError
        });
    }
}
function raSuccess() {
    ra.intro.hide();
    ra.success.show();
}
function raError() {
    ra.intro.hide();
    ra.error.show();
}
function raClose() {
    ra.base.hide();
    ra.intro.show();
    ra.error.hide();
    ra.success.hide();
}
function raOpen() {
    ra.base.show();
}

$(document).ready(function() {
    ra.base = $('#request_access');
    ra.form = $('#request_access_form');
    ra.intro = $('#request_access_intro');
    ra.success = $('#request_access_success');
    ra.error = $('#request_access_error');

    ra.form.on('submit', raSend);
});

$(document).on('keydown', function(event) {
    event.key == "Escape" && raClose();
});

// Clipbord copy functionality
var action_toast_timeout;
function showActionToast(text) {
  clearTimeout(action_toast_timeout);
  var action_toast = $('.action-toast');
  action_toast.text(text).fadeIn()
  action_toast_timeout = setTimeout(function(){ action_toast.fadeOut() }, 5000);
}

$(document).ready(function(){
  new ClipboardJS('[data-snippetcut-btn-name-ru]', {
    text: function(trigger) {
      showActionToast('Скопировано в буфер обмена')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-name]').text();
    }
  });
  new ClipboardJS('[data-snippetcut-btn-name-en]', {
    text: function(trigger) {
      showActionToast('Has been copied to clipboard')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-name]').text();
    }
  });
  new ClipboardJS('[data-snippetcut-btn-text-en]', {
    text: function(trigger) {
      showActionToast('Has been copied to clipboard')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-text]').text();
    }
  });
  new ClipboardJS('[data-snippetcut-btn-text-ru]', {
    text: function(trigger) {
      showActionToast('Скопировано в буфер обмена')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-text]').text();
    }
  });

});

// GDPR

$(document).ready(function(){
    var $gdpr = $('.gdpr');
    var $gdpr_button = $('.gdpr__button');
    var gdpr_status = $.cookie('gdpr-status');

    if (!gdpr_status || gdpr_status != 'accepted') {
        $gdpr.show();
    }

    $gdpr_button.on('click', function() {
        $gdpr.hide();
        $.cookie('gdpr-status', 'accepted', {path: '/' ,  expires: 3650 });
    })
});

//Fixed sidebar
// $(document).ready(function() {
window.onload = function() {
  const wrapper = $('.content');
  const sidebarWrapperInner = $('.sidebar__wrapper-inner');
  const sidebar = $('.sidebar__container');
  const sidebarOffsetTop = sidebar.offset().top - 100;
  const sidebarTopClass = 'sidebar-fixed__top';
  const sidebarBottomClass = 'sidebar-fixed__bottom';
  const footerHeight = $('.footer').height();
  const headerHeight = $('.header').height();
  const docHeight = $(document).height();
  const screenHeight = $(window).outerHeight();

  // console.log($(window).scrollTop(), 'scrollTop');
  // console.log(window.scrollY, 'window.scrollY');
  // console.log($(document).height(), 'document.height');
  // console.log($(window).outerHeight(), 'window.outerHeight');
  // console.log($(document).height() - (319 + $(window).outerHeight()), 'sum');

  // if ($(window).scrollTop() > ($(document).height() - (319 + $(window).outerHeight()))) {
  //   console.log('class bottom');
  //   console.log($(window).scrollTop());
  //   console.log(($(document).height() - (319 + $(window).outerHeight())));
  //
  //
  //   wrapper.removeClass(sidebarTopClass);
  //   wrapper.addClass(sidebarBottomClass);
  //   sidebarWrapperInner.css({
  //     height: $('.layout-sidebar__content').css('height')
  //   });
  //   console.log('in if logs')
  //   console.log($(window).height(), 'window).height');
  //   console.log(headerHeight, 'headerHeight');
  //   console.log(docHeight, 'docHeight');
  //   console.log($(window).scrollTop(), 'window).scrollTop');
  //   console.log(footerHeight, 'footerHeight');
  //   sidebar.css({
  //     height: `calc(100vh - (190px + ${$(window).height() - headerHeight - (docHeight - $(window).scrollTop() - footerHeight)}px))`,
  //     bottom: '25px'
  //   });
  // } else if (wrapper.hasClass(sidebarBottomClass) && $(window).scrollTop() < ($(document).height() - (319 + $(window).outerHeight()))) {
  //   console.log('no class bottom');
  //   wrapper.removeClass(sidebarBottomClass);
  //   wrapper.addClass(sidebarTopClass);
  //   sidebarWrapperInner.css({
  //     height: 'calc(100vh - 130px)'
  //   })
  //   sidebar.css({
  //     height: '',
  //     bottom: 'auto'
  //   });
  // }
  let bottomFixPoint = docHeight - (footerHeight + screenHeight);

  if ($('.layout-sidebar__content').height() > screenHeight) {
    console.log(123)
    setFooterClass(screenHeight, bottomFixPoint, sidebar, wrapper, sidebarTopClass, sidebarBottomClass, sidebarWrapperInner, headerHeight, footerHeight, docHeight)
  } else {
    console.log(3333)
  }


  setTopClass($(window).scrollTop(), sidebarOffsetTop, wrapper, sidebarTopClass, sidebarBottomClass, sidebarWrapperInner, bottomFixPoint);

  $(window).scroll(function() {
    const scrolled = $(this).scrollTop();
    // const footerHeightCalc = docHeight - scrolled - footerHeight;
    bottomFixPoint = $(document).height() - (319 + $(window).outerHeight());


    setTopClass(scrolled, sidebarOffsetTop, wrapper, sidebarTopClass, sidebarBottomClass, sidebarWrapperInner, bottomFixPoint);

    // if (scrolled > sidebarOffsetTop) {
    //   wrapper.addClass(sidebarTopClass);
    // } else if (scrolled < sidebarOffsetTop) {
    //   wrapper.removeClass(sidebarTopClass);
    //   sidebarWrapperInner.css({
    //     height: `calc(100vh - (190px - ${scrolled}px))`
    //   });
    // }

    // const bottomFixPoint = $(document).height() - (sidebarHeight + footerHeight + $(window).outerHeight());

    if ($('.layout-sidebar__content').height() > screenHeight) {
      setFooterClass(scrolled, bottomFixPoint, sidebar, wrapper, sidebarTopClass, sidebarBottomClass, sidebarWrapperInner, headerHeight, footerHeight, docHeight);
    }

    // if (scrolled > bottomFixPoint) {
    //   wrapper.removeClass(sidebarTopClass);
    //   wrapper.addClass(sidebarBottomClass);
    //   sidebarWrapperInner.css({
    //     height: $('.layout-sidebar__content').css('height')
    //   })
    //   sidebar.css({
    //     height: `calc(100vh - (190px + ${$(window).height() - headerHeight - footerHeightCalc}px))`,
    //     bottom: '25px'
    //   });
    // } else if (wrapper.hasClass(sidebarBottomClass) && scrolled < bottomFixPoint) {
    //   wrapper.removeClass(sidebarBottomClass);
    //   wrapper.addClass(sidebarTopClass);
    //   sidebarWrapperInner.css({
    //     height: 'calc(100vh - 130px)'
    //   })
    //   sidebar.css({
    //     height: '',
    //     bottom: 'auto'
    //   });
    // }
  });
};

function setTopClass(scrolled, sidebarOffsetTop, wrapper, classTop, classBottom, sidebarWrapper, bottomFixPoint) {
  if (scrolled > sidebarOffsetTop && scrolled < bottomFixPoint) {
    wrapper.addClass(classTop);
    sidebarWrapper.css({
      height: `calc(100vh - 130px)`
    });
  } else if (scrolled < sidebarOffsetTop || wrapper.hasClass(classBottom)) {
    wrapper.removeClass(classTop);
    sidebarWrapper.css({
      height: `calc(100vh - (190px - ${scrolled}px))`
    });
  }
}

function setFooterClass(scrolled, bottomFixPoint, sidebar, wrapper, classTop, classBottom, sidebarWrapper, headerHeight, footerHeight, docHeight) {
  const footerHeightCalc = docHeight - scrolled - footerHeight;

  if (scrolled > bottomFixPoint) {
    wrapper.removeClass(classTop);
    wrapper.addClass(classBottom);
    sidebarWrapper.css({
      height: $('.layout-sidebar__content').css('height')
    })
    sidebar.css({
      height: `calc(100vh - (190px + ${$(window).height() - headerHeight - footerHeightCalc}px))`,
      bottom: '25px'
    });
  } else if (wrapper.hasClass(classBottom) && scrolled < bottomFixPoint) {
    wrapper.removeClass(classBottom);
    wrapper.addClass(classTop);
    sidebarWrapper.css({
      height: 'calc(100vh - 130px)'
    })
    sidebar.css({
      height: '',
      bottom: 'auto'
    });
  }
}
