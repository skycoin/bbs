import {Component, HostBinding, OnInit, ViewEncapsulation} from '@angular/core';

@Component({
    selector: 'app-loading',
    templateUrl: './loading.component.html',
    styleUrls: ['./loading.component.scss'],
    encapsulation: ViewEncapsulation.None
})
export class LoadingComponent implements OnInit {
    loadingText = 'Loading';
    @HostBinding('style.display') display = 'none';

    constructor() {
    }

    ngOnInit() {
    }

    /**
     * Show Loading And Return Promise
     * @param loadingText Loading Text
     */
    start(loadingText: string = 'Loading') {
        this.loadingText = loadingText;
        this.display = 'block';
        return Promise.resolve();
    }

    close() {
        this.loadingText = '';
        this.display = 'none';
    }
}
