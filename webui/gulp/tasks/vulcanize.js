'use strict';

// Vulcanize imports
module.exports = function (gulp, plugins, config) { return function () {
  return gulp.src('dist/elements/elements.vulcanized.html')
    .pipe(plugins.vulcanize({
      dest: 'dist/elements',
      inlineScripts: false,
    }))
    .pipe(plugins.crisper())
    .pipe(gulp.dest('dist/elements'))
    .pipe(plugins.size({title: 'vulcanize'}));
};};
