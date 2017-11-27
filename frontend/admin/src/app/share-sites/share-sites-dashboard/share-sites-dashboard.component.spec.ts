import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShareSitesDashboardComponent } from './share-sites-dashboard.component';

describe('ShareSitesDashboardComponent', () => {
  let component: ShareSitesDashboardComponent;
  let fixture: ComponentFixture<ShareSitesDashboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShareSitesDashboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShareSitesDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
