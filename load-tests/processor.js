// Artillery processor functions for dynamic data generation

module.exports = {
  // Generate random product type
  generateProductType: function(context, events, done) {
    const types = ['server', 'client'];
    context.vars.productType = types[Math.floor(Math.random() * types.length)];
    return done();
  },

  // Generate random version number
  generateVersionNumber: function(context, events, done) {
    const major = Math.floor(Math.random() * 10) + 1;
    const minor = Math.floor(Math.random() * 10);
    const patch = Math.floor(Math.random() * 10);
    context.vars.versionNumber = `${major}.${minor}.${patch}`;
    return done();
  },

  // Generate random release type
  generateReleaseType: function(context, events, done) {
    const types = ['major', 'minor', 'patch', 'feature', 'security'];
    context.vars.releaseType = types[Math.floor(Math.random() * types.length)];
    return done();
  },

  // Generate random notification type
  generateNotificationType: function(context, events, done) {
    const types = ['info', 'warning', 'error', 'success'];
    context.vars.notificationType = types[Math.floor(Math.random() * types.length)];
    return done();
  },

  // Generate unique product ID
  generateProductId: function(context, events, done) {
    const timestamp = Date.now();
    const random = Math.random().toString(36).substring(2, 8);
    context.vars.productId = `load-test-${random}-${timestamp}`;
    return done();
  }
};

