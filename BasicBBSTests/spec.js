// spec.js
describe('BBS page loaded - Create identity sub screen appears', function() {
  var title ='Skycoin BBS'
  var loginBtn = element(by.className('btn close-btn'));
  var exitBtn = element(by.className('btn create-btn'));
  var createIdnLnk = element(by.className('ng-tns-c0-0'));
  var createIdnImg = element(by.className('fa fa-plus-circle'));//?not sure locator is accurate at all
  var IdentityScrTitle = element(by.className('modal-title'));
  
  beforeEach(function(){
    browser.get('http://localhost:8080/');
  });

  it('Title is correct , page loaded', function() {
    expect(browser.getTitle()).
    toEqual(title); 
  })
  it('Identity Screen , exit button appears', function() {

    expect(exitBtn.isPresent()).
        toEqual(true); 
  })
  it('Identity Screen , login button appears', function() {

    expect(loginBtn.isPresent()).
        toEqual(true); 
  })
  it('Identity Screen , link appears', function() {
    expect(loginBtn.isPresent()).
        toEqual(true);
  })
})