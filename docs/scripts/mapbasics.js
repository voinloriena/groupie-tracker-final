ymaps.ready(function () {
    var MyMap = new ymaps.Map('map', {
        center: [51.15296161002877,71.48300799999997],
        zoom: 5,
        searchControlProvider: 'yandex#search'
    });
    var alem = new ymaps.Placemark([51.088352098944235,71.41420656258758]);
    MyMap.geoObjects.add(alem);

    var myCollection = new ymaps.GeoObjectCollection();
    var npm = document.getElementsByClassName('placemark');
    myGeocoder = ymaps.geocode("astana");
    for (var i = 0; i < npm.length; ++i) {
    // do something with npm[i], which is a <li> element
    
        myGeocoder = ymaps.geocode(npm[i].textContent.replace(/[-_]/g, " "));
            console.log(npm[i].textContent.replace(/[-_]/g, " "))
            myGeocoder.then(
                function (res) {
                    MyMap.geoObjects.add(res.geoObjects.get(1));
                    // Выведем в консоль данные, полученные в результате геокодирования объекта.
                    console.log(res.geoObjects.get(1).properties.get('metaDataProperty').getAll());
                },
                function (err) {
                    // Обработка ошибки.
                }
            );
    }
    MyMap.geoObjects.add(myCollection);

});
