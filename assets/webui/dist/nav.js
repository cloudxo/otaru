import {$, $$} from './domhelper.js';

const reContentHash = /^#([\w-]+)$/;

const validContents = [];
$$(".nav__tab").forEach(tab => {
  const contentId = tab.getAttribute('href').match(reContentHash)[1];
  validContents.push(contentId);
});

let currContentId = '';
let currBrowsefsPath = '//';

const updateContentIfNeeded = (opts) => {
  let refreshNeeded = false;
  const fromHash = (opts['fromHash'] === true);

  let newContentId = opts['contentId'];
  let newBrowsefsPath = opts['currBrowsefsPath'];
  if (fromHash) {
    let m = window.location.hash.match(/^#(\w+)(\/.*)?$/);
    if (m) {
      newContentId = m[1];
      newBrowsefsPath = m[2] ? decodeURI(m[2]) : undefined;
    }
  }

  if (newContentId !== undefined && newContentId !== currContentId) {
    currContentId = newContentId;
    refreshNeeded = true;
  }
  if (newBrowsefsPath !== undefined && newBrowsefsPath !== currBrowsefsPath) {
    currBrowsefsPath = newBrowsefsPath;
    refreshNeeded = true;
  }

  if (refreshNeeded) {
    const pushHistory = !fromHash;
    showContent(pushHistory);
  }
};

$$(".nav__tab").forEach(tabA => {
  tabA.addEventListener("click", ev => {
    ev.preventDefault();

    const newContentId = tabA.getAttribute('href').match(reContentHash)[1];
    
    updateContentIfNeeded({contentId: newContentId});
    return false; 
  });
});

const contentSection = contentId => {
  if (!validContents.includes(contentId))
    throw new Error(`Invalid contentId "${contentId}"`);

  return $(`.section--${contentId}`)
};
const isSectionSelected = contentId =>
  contentSection(contentId).classList.contains('section--selected');

const showContent = (pushHistory = false) => {
  let newHref = `#${currContentId}`;
  if (currContentId === 'browsefs') {
    newHref += encodeURIComponent(currBrowsefsPath).replace(/%2f/gi, '/');
  }
  const newTitle = `Otaru WebUI: ${currContentId}`;

  const modState = pushHistory ? window.history.pushState : window.history.replaceState;
  document.title = newTitle;
  modState.call(window.history, null, newTitle, newHref);

  if (!validContents.includes(currContentId)) {
    console.log(`invalid contentId "${currContentId}"`);
    return;
  }

  const contentHash = `#${currContentId}`;
  $$(".nav__tab").forEach(tab => {
    if (tab.getAttribute("href") === contentHash) {
      tab.classList.add("nav__item--selected");
    } else {
      tab.classList.remove("nav__item--selected");
    }
  });
  const sectionNeedle = `section--${currContentId}`;

  $$(".section").forEach(section => {
    if (section.classList.contains(sectionNeedle)) {
      section.classList.add("section--selected");
      const e = new Event('shown');
      e.browsefsPath = currBrowsefsPath;
      section.dispatchEvent(e);
    } else {
      section.classList.remove("section--selected");
      section.dispatchEvent(new Event('hidden'));
    }
  });
};

window.addEventListener("hashchange", () => {
  updateContentIfNeeded({fromHash: true});
})
window.addEventListener("DOMContentLoaded", () => {
  updateContentIfNeeded({fromHash: true});
}, {oneshot: true});

export {contentSection, isSectionSelected, updateContentIfNeeded};
