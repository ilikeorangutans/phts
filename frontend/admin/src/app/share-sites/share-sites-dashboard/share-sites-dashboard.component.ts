import { ShareSiteService } from './../../services/share-site.service';
import { Component, OnInit } from '@angular/core';
import { ShareSite } from '../../models/share-site';

@Component({
  selector: 'app-share-sites-dashboard',
  templateUrl: './share-sites-dashboard.component.html',
  styleUrls: ['./share-sites-dashboard.component.css']
})
export class ShareSitesDashboardComponent implements OnInit {

  shareSites: Array<ShareSite> = new Array();

  constructor(
    private shareSiteService: ShareSiteService
  ) { }

  ngOnInit() {
    this.shareSiteService.list().then(sites => {
      this.shareSites = sites;
    });
  }

}
