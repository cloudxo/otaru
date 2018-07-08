import {contentSection, updateContentIfNeeded} from './nav.js';
import {fillRemoteContent} from './api.js';
import {$} from './domhelper.js';
import './browsefs.js';

const leftfs = $("#leftfs");
const rightfs = $("#rightfs");

contentSection('browsefs').addEventListener('shown', e => {
  if (e.browsefsPath !== undefined)
    rightfs.path = e.browsefsPath;
  rightfs.triggerUpdate();
});
contentSection('browsefs').addEventListener('hidden', () => {
  return rightfs.clear();
});
rightfs.addEventListener('pathChanged', e => {
  updateContentIfNeeded({currBrowsefsPath: e.newPath});
});

const splitbar = $('.splitbar');
const noophandler = () => { return false; };
splitbar.addEventListener('mousedown', md => {
  const pn = splitbar.parentNode;
  const offX = pn.offsetLeft;
  const offW = pn.offsetWidth;
  const leftpane = pn.querySelector('.split--leftpane');
  const rightpane = pn.querySelector('.split--rightpane');

  const mmhandler = mm => {
    const l = (event.pageX - offX) / offW;
    leftpane.style.width = `${l * 100}%`;
    rightpane.style.width =`${(1.0 - l) * 100}%`;
  };
  const muhandler = mu => {
    pn.removeEventListener('mousemove', mmhandler);
    pn.removeEventListener('mouseup', muhandler);
    pn.removeEventListener('selectstart', noophandler);
    pn.classList.remove('drag_parent');
  };
  pn.addEventListener('mousemove', mmhandler);
  pn.addEventListener('mouseup', muhandler);
  pn.addEventListener('selectstart', noophandler);
  pn.classList.add('drag_parent');
});

window.addEventListener('DOMContentLoaded', () => {
  let focusfs = leftfs;
  focusfs.cursorIndex = 0;

  document.addEventListener('keypress', e => {
    if (e.target instanceof HTMLInputElement)
      return true;

    if (e.key === 'j') {
      focusfs.setCursorIndexBounded(focusfs.cursorIndex + 1);
    } else if (e.key === 'k') {
      focusfs.setCursorIndexBounded(Math.max(focusfs.cursorIndex - 1, 0));
    } else if (e.key === 'l') {
      rightfs.setCursorIndexBounded(leftfs.cursorIndex);
      leftfs.clearCursor();
      focusfs = rightfs;
    } else if (e.key === 'h') {
      leftfs.setCursorIndexBounded(rightfs.cursorIndex);
      rightfs.clearCursor();
      focusfs = leftfs;
    } else if (e.key === ' ') {
      let cr = focusfs.cursorRow;
      if (cr)
        cr.toggleSelection(); 
    } else if (e.key === 'Enter') {
      let cr = focusfs.cursorRow;
      if (cr)
        cr.triggerAction(); 
    } else if (e.key === 'u') {
      focusfs.navigateParent();
    } else if (e.key === '?') {
      focusfs.query = '';
    } else {
      console.log(`keypress ${e.key}`);
    }
  });
  document.addEventListener('keyup', e => {
    if (e.key === 'Escape') {
      focusfs.query = null; 
    }
    return false;
  });
});
