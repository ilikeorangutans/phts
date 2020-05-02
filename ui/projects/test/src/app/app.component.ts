import { Component } from '@angular/core';

import { VersionServiceClient } from '../proto/phts.pbsc';
import { VersionRequest, VersionResponse } from '../proto/phts.pb';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'test';

  constructor(
    readonly v: VersionServiceClient
  ) {


  }

  click() {
    console.log('click');
    this.v.get(new VersionRequest()).subscribe(
      response => console.log(response.toJSON),
      error => console.log("got grpc error", error),
    );

  }
}
