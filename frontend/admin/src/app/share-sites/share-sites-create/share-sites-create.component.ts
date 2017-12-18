import { ShareSiteService } from './../../services/share-site.service';
import { ShareSite } from './../../models/share-site';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-share-sites-create',
  templateUrl: './share-sites-create.component.html',
  styleUrls: ['./share-sites-create.component.css']
})
export class ShareSitesCreateComponent implements OnInit {

  shareSite = new ShareSite();

  constructor(
    private shareSiteService: ShareSiteService
  ) { }

  ngOnInit() {
  }

  onSubmit() {
    console.log('onSubmit()');
    this.shareSiteService.save(this.shareSite).then(shareSite => console.log(shareSite)).catch(e => console.log(e));
  }
}
