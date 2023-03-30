var data = [];

// Remove and complete icons in SVG format
var removeSVG =
  '<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 22 22" style="enable-background:new 0 0 22 22;" xml:space="preserve"><rect class="noFill" width="22" height="22"/><g><g><path class="fill" d="M16.1,3.6h-1.9V3.3c0-1.3-1-2.3-2.3-2.3h-1.7C8.9,1,7.8,2,7.8,3.3v0.2H5.9c-1.3,0-2.3,1-2.3,2.3v1.3c0,0.5,0.4,0.9,0.9,1v10.5c0,1.3,1,2.3,2.3,2.3h8.5c1.3,0,2.3-1,2.3-2.3V8.2c0.5-0.1,0.9-0.5,0.9-1V5.9C18.4,4.6,17.4,3.6,16.1,3.6z M9.1,3.3c0-0.6,0.5-1.1,1.1-1.1h1.7c0.6,0,1.1,0.5,1.1,1.1v0.2H9.1V3.3z M16.3,18.7c0,0.6-0.5,1.1-1.1,1.1H6.7c-0.6,0-1.1-0.5-1.1-1.1V8.2h10.6V18.7z M17.2,7H4.8V5.9c0-0.6,0.5-1.1,1.1-1.1h10.2c0.6,0,1.1,0.5,1.1,1.1V7z"/></g><g><g><path class="fill" d="M11,18c-0.4,0-0.6-0.3-0.6-0.6v-6.8c0-0.4,0.3-0.6,0.6-0.6s0.6,0.3,0.6,0.6v6.8C11.6,17.7,11.4,18,11,18z"/></g><g><path class="fill" d="M8,18c-0.4,0-0.6-0.3-0.6-0.6v-6.8c0-0.4,0.3-0.6,0.6-0.6c0.4,0,0.6,0.3,0.6,0.6v6.8C8.7,17.7,8.4,18,8,18z"/></g><g><path class="fill" d="M14,18c-0.4,0-0.6-0.3-0.6-0.6v-6.8c0-0.4,0.3-0.6,0.6-0.6c0.4,0,0.6,0.3,0.6,0.6v6.8C14.6,17.7,14.3,18,14,18z"/></g></g></g></svg>';

renderChecks();

function renderChecks() {
fetch("/checks",{
    method: "get",
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json"
    },
  }).then(response => {
    if (response.status >= 200 || response.status < 500) {
      response.json()
      .then(json => {
        json.forEach(order => {
          console.log(`adding ${order.orderId} - ${order.total }`);
          addItemToDOM(order.orderId, order.total, order.items, order.url);
        })
      })
    } else {
      console.log(`error ${response.status}`)  
    }
  }).catch(error => { 
    console.log(error)
  })
}

function removeItem() {
  var item = this.parentNode.parentNode;
  var parent = item.parentNode;
  var orderId = item.id;

  fetch(`/checks/${orderId}`, {
    method: "DELETE",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  })
    .then((response) => {
      if (response.status >= 200 || response.status < 500) {
        parent.removeChild(item);
      } else {
        console.log(`error ${response.status}`);
      }
    })
    .catch((error) => {
      console.log(error);
    });
}


// Adds a new item to the fod list
function addItemToDOM(checkId, total, items, url) {
  var list = document.getElementById("check");
  var item = document.createElement("li");
  const food = items.map(i => i.name).join(", ");
  
  item.innerHTML = `<div>Check Id ${checkId} - Total $${total}: ${food} - <a href="${url}">Download Receipt</a></div>`;
  item.id = checkId;

  var buttons = document.createElement("div");
  buttons.classList.add("buttons");

  var remove = document.createElement("button");
  remove.classList.add("remove");
  remove.innerHTML = removeSVG;

  // Add click event for removing the item
  remove.addEventListener("click", removeItem);

  buttons.appendChild(remove);
  item.appendChild(buttons);

  list.insertBefore(item, list.childNodes[0]);
}