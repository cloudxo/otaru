import {contentSection, isSectionSelected} from './nav.js';
import {rpc, fillRemoteContent} from './api.js';
import {$, $$, removeAllChildNodes} from './util.js';

const updateInterval = 3000;

(() => {
  const colNames = ['type', 'name', 'size', 'uid', 'gid', 'perm_mode', 'modified_time'];
  const reTime = /_time$/;
  const sortFuncMap = {
    'name': (a, b) => {
      if (a['type'] != b['type']) {
        return a['type'].charCodeAt() - b['type'].charCodeAt();
      }
      return a['name'].localeCompare(b['name']);
    },
    'time_asc': (a, b) => a['modified_time'] - b['modified_time'],
    'time_desc': (a, b) => b['modified_time'] - a['modified_time'],
  };

  const triggerUpdate = async () => {
    if (!isSectionSelected('browsefs'))
      return;

    try {
      const result = await rpc('/v1/filesystem/ls?path=/');
      const listDiv = $('.browsefs__list');
      removeAllChildNodes($('.browsefs__list'));

      const sortSel = $('.browsefs__sort').value;
      const sortFunc = sortFuncMap[sortSel];

      if (result['entry'] === undefined) {
        listDiv.classList.add('.browsefs__list--empty');
      } else {
        const rows = result['entry'].sort(sortFunc);
        for (let row of rows) {
          listDiv.classList.remove('.browsefs__list--empty');

          var rowDiv = document.createElement('tr');
          rowDiv.classList.add('browsefs__entry');
          listDiv.appendChild(rowDiv);

          for (let colName of colNames) {
            var cell = document.createElement('td');
            cell.classList.add('browsefs__cell');
            cell.classList.add(`browsefs__cell--${colName}`);

            let val = row[colName];
            if (colName === 'type') {
              val = val[0];
            } else if (colName === 'perm_mode') {
              val = val.toString(8);
            } else if (colName === 'uid') {
              val = 'u'+val;
            } else if (colName === 'gid') {
              val = 'g'+val;
            } else if (colName === 'size') {
              if (val === undefined) {
                val = '';
              } else if (val > 1024 * 1024 * 1024) {
                val = (val / (1024 * 1024 * 1024)).toPrecision(2) + 'GiB';
              } else if (val > 1024 * 1024) {
                val = (val / (1024 * 1024)).toPrecision(2) + 'MiB';
              } else if (val > 1024) {
                val = (val / 1024).toPrecision(2) + 'KiB';
              } else {
                val = val + 'B';
              }
            } else if (colName.match(reTime)) {
              const t = new Date(val*1000);
              const diff = new Date() - t;

              const startOfToday = new Date();
              startOfToday.setHours(0);
              startOfToday.setMinutes(0);
              startOfToday.setSeconds(0);
              startOfToday.setMilliseconds(0);

              const pad = n => (n < 10 ? '0' : '') + n;
              if (diff < 60 * 1000) {
                val = `${(diff / (1000)).toFixed(0)}s ago`;
              } else if (diff < 1 * 60 * 60 * 1000) {
                val = `${(diff / (60 * 1000)).toFixed(0)}m ago`;
              } else if (diff < 6 * 60 * 60 * 1000) {
                val = `${(diff / (60 * 60 * 1000)).toFixed(0)}h ago`;
              } else if (t > startOfToday) {
                val = `${pad(t.getHours())}:${pad(t.getMinutes())}`;
              } else {
                val = `${pad(t.getFullYear()-2000)}/${pad(t.getMonth()+1)}/${pad(t.getDay()+1)}`;
              }
            }
            if (val === undefined)
              val = '-';

            cell.textContent = val;

            rowDiv.appendChild(cell);
          }
        }
      }
    } catch (e) {
      console.log(e);
    }
  }
  contentSection('browsefs').addEventListener('shown', triggerUpdate);
  $('.browsefs__sort').addEventListener('change', triggerUpdate);
  contentSection('browsefs').addEventListener('hidden', e => {
    // $('.browsefs__entry').removeChild();
  });
})();

(() => {
  const triggerUpdate = async () => {
    if (!isSectionSelected('blobstore'))
      return;

    try {
      await fillRemoteContent('/v1/blobstore/config', '#blobstore-', [
          'backend_impl_name', 'backend_flags',
          'cache_impl_name', 'cache_flags']);
    } catch (e) {
      console.log(e);
    }
    window.setTimeout(triggerUpdate, updateInterval);
  }
  contentSection('blobstore').addEventListener('shown', e => {
    triggerUpdate();
  });
})();

(() => {
  const triggerUpdate = async () => {
    if (!isSectionSelected('settings'))
      return;

    try {
      await fillRemoteContent("/v1/system/info", "#settings-", [
          'go_version', 'os', 'arch', 'num_goroutine', 'hostname', 'pid', 'uid',
          'mem_alloc', 'mem_sys', 'num_gc', 'num_fds']);
    } catch (e) {
      console.log(e);
    }
    window.setTimeout(triggerUpdate, updateInterval);
  }
  contentSection('settings').addEventListener('shown', e => {
    triggerUpdate();
  });
})();
