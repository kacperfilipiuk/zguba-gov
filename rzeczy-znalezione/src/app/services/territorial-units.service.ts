// territorial-units.service.ts
import { Injectable } from '@angular/core';
import territorialUnitsData from '../../assets/territorial-units.json';

export interface TerritorialUnit {
  name: string;
  type: 'wojewodztwo' | 'powiat' | 'gmina' | 'miasto';
  parentName?: string; // np. dla powiatu: nazwa województwa
  fullName: string; // pełna nazwa z kontekstem
  id?: string; // ID z pliku
  email?: string; // Email z pliku
  officeName?: string; // Nazwa urzędu
  voivodeship?: string; // Województwo
  county?: string; // Powiat
}

interface LocalDataUnit {
  id: string;
  name: string;
  parentId: string;
  level: number;
  kind: string;
  type: 'wojewodztwo' | 'powiat' | 'gmina' | 'miasto';
  hasDescription: boolean;
  email: string;
  officeName: string;
  voivodeship: string;
  county: string;
}

@Injectable({
  providedIn: 'root'
})
export class TerritorialUnitsService {
  private units: TerritorialUnit[] = [];
  private wojewodztwaMap = new Map<string, string>(); // id -> nazwa
  private loading = false;
  private loaded = false;

  constructor() {
    this.loadDataFromLocalFile();
  }

  private async loadDataFromLocalFile(): Promise<void> {
    if (this.loading || this.loaded) return;
    
    this.loading = true;
    
    try {
      console.log('🔄 Rozpoczynam ładowanie danych z lokalnego pliku...');
      
      const data = territorialUnitsData as LocalDataUnit[];
      
      // Mapuj dane z pliku do struktury używanej w aplikacji
      this.units = data.map(unit => {
        let fullName = unit.name;
        
        // Dodaj kontekst dla pełnej nazwy
        if (unit.type === 'gmina' || unit.type === 'powiat') {
          fullName = `${unit.name} (woj. ${unit.voivodeship})`;
        }
        
        return {
          id: unit.id,
          name: unit.name,
          type: unit.type,
          parentName: unit.voivodeship,
          fullName: fullName,
          email: unit.email,
          officeName: unit.officeName,
          voivodeship: unit.voivodeship,
          county: unit.county
        };
      });
      
      this.loaded = true;
      console.log(`✅ Załadowano łącznie ${this.units.length} jednostek terytorialnych z lokalnego pliku`);
      console.log(`   - Województwa: ${this.units.filter(u => u.type === 'wojewodztwo').length}`);
      console.log(`   - Powiaty: ${this.units.filter(u => u.type === 'powiat').length}`);
      console.log(`   - Miasta: ${this.units.filter(u => u.type === 'miasto').length}`);
      console.log(`   - Gminy: ${this.units.filter(u => u.type === 'gmina').length}`);
    } catch (error) {
      console.error('❌ Błąd podczas ładowania danych z lokalnego pliku:', error);
      this.loadBackupData();
    } finally {
      this.loading = false;
    }
  }

  private async loadWojewodztwa(): Promise<void> {
    // Ta metoda nie jest już używana - dane są ładowane z pliku lokalnego
  }

  private async loadPowiatyIMiasta(): Promise<void> {
    // Ta metoda nie jest już używana - dane są ładowane z pliku lokalnego
  }

  private async loadGminy(): Promise<void> {
    // Ta metoda nie jest już używana - dane są ładowane z pliku lokalnego
  }

  private getWojewodztwoFromParentId(parentId: string): string {
    // Ta metoda nie jest już używana
    return 'nieznane';
  }

  private getWojewodztwoFromUnitId(unitId: string): string {
    // Ta metoda nie jest już używana
    return 'nieznane';
  }

  private formatWojewodztwoName(name: string): string {
    // Ta metoda nie jest już używana
    return name;
  }

  private loadBackupData(): void {
    // Województwa
    const wojewodztwa = [
      'dolnośląskie', 'kujawsko-pomorskie', 'lubelskie', 'lubuskie',
      'łódzkie', 'małopolskie', 'mazowieckie', 'opolskie',
      'podkarpackie', 'podlaskie', 'pomorskie', 'śląskie',
      'świętokrzyskie', 'warmińsko-mazurskie', 'wielkopolskie', 'zachodniopomorskie'
    ];

    wojewodztwa.forEach(woj => {
      this.units.push({
        name: `Województwo ${woj}`,
        type: 'wojewodztwo',
        fullName: `Województwo ${woj}`
      });
    });

    // Przykładowe powiaty (reprezentatywna lista najważniejszych)
    const powiaty = [
      { name: 'Powiat warszawski zachodni', wojewodztwo: 'mazowieckie' },
      { name: 'Powiat piaseczyński', wojewodztwo: 'mazowieckie' },
      { name: 'Powiat pruszkowski', wojewodztwo: 'mazowieckie' },
      { name: 'Powiat wołomiński', wojewodztwo: 'mazowieckie' },
      { name: 'Powiat grodziski', wojewodztwo: 'mazowieckie' },
      { name: 'Powiat krakowski', wojewodztwo: 'małopolskie' },
      { name: 'Powiat wielicki', wojewodztwo: 'małopolskie' },
      { name: 'Powiat myślenicki', wojewodztwo: 'małopolskie' },
      { name: 'Powiat poznański', wojewodztwo: 'wielkopolskie' },
      { name: 'Powiat wrocławski', wojewodztwo: 'dolnośląskie' },
      { name: 'Powiat gdański', wojewodztwo: 'pomorskie' },
      { name: 'Powiat katowicki', wojewodztwo: 'śląskie' },
      { name: 'Powiat bielski', wojewodztwo: 'śląskie' },
      { name: 'Powiat gliwicki', wojewodztwo: 'śląskie' },
      { name: 'Powiat lubliniecki', wojewodztwo: 'śląskie' },
      { name: 'Powiat łódzki wschodni', wojewodztwo: 'łódzkie' },
      { name: 'Powiat rzeszowski', wojewodztwo: 'podkarpackie' },
      { name: 'Powiat białostocki', wojewodztwo: 'podlaskie' },
      { name: 'Powiat szczeciński', wojewodztwo: 'zachodniopomorskie' },
    ];

    powiaty.forEach(powiat => {
      this.units.push({
        name: powiat.name,
        type: 'powiat',
        parentName: powiat.wojewodztwo,
        fullName: `${powiat.name} (woj. ${powiat.wojewodztwo})`
      });
    });

    // Miasta na prawach powiatu
    const miastaNaPrawachPowiatu = [
      { name: 'Warszawa', wojewodztwo: 'mazowieckie' },
      { name: 'Kraków', wojewodztwo: 'małopolskie' },
      { name: 'Wrocław', wojewodztwo: 'dolnośląskie' },
      { name: 'Poznań', wojewodztwo: 'wielkopolskie' },
      { name: 'Gdańsk', wojewodztwo: 'pomorskie' },
      { name: 'Szczecin', wojewodztwo: 'zachodniopomorskie' },
      { name: 'Bydgoszcz', wojewodztwo: 'kujawsko-pomorskie' },
      { name: 'Lublin', wojewodztwo: 'lubelskie' },
      { name: 'Katowice', wojewodztwo: 'śląskie' },
      { name: 'Białystok', wojewodztwo: 'podlaskie' },
      { name: 'Gdynia', wojewodztwo: 'pomorskie' },
      { name: 'Częstochowa', wojewodztwo: 'śląskie' },
      { name: 'Radom', wojewodztwo: 'mazowieckie' },
      { name: 'Sosnowiec', wojewodztwo: 'śląskie' },
      { name: 'Toruń', wojewodztwo: 'kujawsko-pomorskie' },
      { name: 'Kielce', wojewodztwo: 'świętokrzyskie' },
      { name: 'Rzeszów', wojewodztwo: 'podkarpackie' },
      { name: 'Gliwice', wojewodztwo: 'śląskie' },
      { name: 'Zabrze', wojewodztwo: 'śląskie' },
      { name: 'Olsztyn', wojewodztwo: 'warmińsko-mazurskie' },
      { name: 'Bielsko-Biała', wojewodztwo: 'śląskie' },
      { name: 'Bytom', wojewodztwo: 'śląskie' },
      { name: 'Zielona Góra', wojewodztwo: 'lubuskie' },
      { name: 'Rybnik', wojewodztwo: 'śląskie' },
      { name: 'Ruda Śląska', wojewodztwo: 'śląskie' },
      { name: 'Opole', wojewodztwo: 'opolskie' },
      { name: 'Tychy', wojewodztwo: 'śląskie' },
      { name: 'Gorzów Wielkopolski', wojewodztwo: 'lubuskie' },
      { name: 'Elbląg', wojewodztwo: 'warmińsko-mazurskie' },
      { name: 'Dąbrowa Górnicza', wojewodztwo: 'śląskie' },
      { name: 'Płock', wojewodztwo: 'mazowieckie' },
      { name: 'Wałbrzych', wojewodztwo: 'dolnośląskie' },
      { name: 'Włocławek', wojewodztwo: 'kujawsko-pomorskie' },
      { name: 'Tarnów', wojewodztwo: 'małopolskie' },
      { name: 'Chorzów', wojewodztwo: 'śląskie' },
      { name: 'Koszalin', wojewodztwo: 'zachodniopomorskie' },
      { name: 'Kalisz', wojewodztwo: 'wielkopolskie' },
      { name: 'Legnica', wojewodztwo: 'dolnośląskie' },
      { name: 'Grudziądz', wojewodztwo: 'kujawsko-pomorskie' },
      { name: 'Słupsk', wojewodztwo: 'pomorskie' },
      { name: 'Jaworzno', wojewodztwo: 'śląskie' },
      { name: 'Jastrzębie-Zdrój', wojewodztwo: 'śląskie' },
      { name: 'Nowy Sącz', wojewodztwo: 'małopolskie' },
      { name: 'Jelenia Góra', wojewodztwo: 'dolnośląskie' },
      { name: 'Siedlce', wojewodztwo: 'mazowieckie' },
      { name: 'Mysłowice', wojewodztwo: 'śląskie' },
      { name: 'Konin', wojewodztwo: 'wielkopolskie' },
      { name: 'Piotrków Trybunalski', wojewodztwo: 'łódzkie' },
      { name: 'Inowrocław', wojewodztwo: 'kujawsko-pomorskie' },
      { name: 'Lubin', wojewodztwo: 'dolnośląskie' },
      { name: 'Ostrów Wielkopolski', wojewodztwo: 'wielkopolskie' },
      { name: 'Suwałki', wojewodztwo: 'podlaskie' },
      { name: 'Stargard', wojewodztwo: 'zachodniopomorskie' },
      { name: 'Gniezno', wojewodztwo: 'wielkopolskie' },
      { name: 'Głogów', wojewodztwo: 'dolnośląskie' },
      { name: 'Pabianice', wojewodztwo: 'łódzkie' },
      { name: 'Chełm', wojewodztwo: 'lubelskie' },
      { name: 'Zamość', wojewodztwo: 'lubelskie' },
      { name: 'Tomaszów Mazowiecki', wojewodztwo: 'łódzkie' },
      { name: 'Przemyśl', wojewodztwo: 'podkarpackie' },
      { name: 'Stalowa Wola', wojewodztwo: 'podkarpackie' },
      { name: 'Mielec', wojewodztwo: 'podkarpackie' },
      { name: 'Krosno', wojewodztwo: 'podkarpackie' },
      { name: 'Tarnobrzeg', wojewodztwo: 'podkarpackie' },
      { name: 'Łomża', wojewodztwo: 'podlaskie' },
      { name: 'Biała Podlaska', wojewodztwo: 'lubelskie' },
      { name: 'Skierniewice', wojewodztwo: 'łódzkie' },
      { name: 'Sopot', wojewodztwo: 'pomorskie' },
      { name: 'Świnoujście', wojewodztwo: 'zachodniopomorskie' }
    ];

    miastaNaPrawachPowiatu.forEach(miasto => {
      this.units.push({
        name: miasto.name,
        type: 'miasto',
        parentName: miasto.wojewodztwo,
        fullName: `Miasto ${miasto.name} (woj. ${miasto.wojewodztwo})`
      });
    });

    // Przykładowe gminy (reprezentatywne dla największych aglomeracji)
    const gminy = [
      { name: 'Gmina Lesznowola', wojewodztwo: 'mazowieckie' },
      { name: 'Gmina Konstancin-Jeziorna', wojewodztwo: 'mazowieckie' },
      { name: 'Gmina Michałowice', wojewodztwo: 'mazowieckie' },
      { name: 'Gmina Raszyn', wojewodztwo: 'mazowieckie' },
      { name: 'Gmina Wieliczka', wojewodztwo: 'małopolskie' },
      { name: 'Gmina Zabierzów', wojewodztwo: 'małopolskie' },
      { name: 'Gmina Tarnowo Podgórne', wojewodztwo: 'wielkopolskie' },
      { name: 'Gmina Komorniki', wojewodztwo: 'wielkopolskie' },
      { name: 'Gmina Długołęka', wojewodztwo: 'dolnośląskie' },
      { name: 'Gmina Kobierzyce', wojewodztwo: 'dolnośląskie' },
      { name: 'Gmina Żukowo', wojewodztwo: 'pomorskie' },
      { name: 'Gmina Pruszcz Gdański', wojewodztwo: 'pomorskie' },
    ];

    gminy.forEach(gmina => {
      this.units.push({
        name: gmina.name,
        type: 'gmina',
        parentName: gmina.wojewodztwo,
        fullName: `${gmina.name} (woj. ${gmina.wojewodztwo})`
      });
    });
  }

  /**
   * Wyszukuje jednostki terytorialne na podstawie zapytania
   * @param query - fraza do wyszukania
   * @param type - opcjonalny filtr typu jednostki
   * @returns tablica pasujących jednostek
   */
  async search(query: string, type?: TerritorialUnit['type']): Promise<TerritorialUnit[]> {
    // Czekaj na załadowanie danych jeśli jeszcze nie są gotowe
    if (this.loading) {
      await this.waitForLoad();
    }
    
    if (!query || query.length < 2) {
      return [];
    }

    const normalizedQuery = this.normalize(query);
    
    let filtered = this.units.filter(unit => {
      const matchesQuery = this.normalize(unit.fullName).includes(normalizedQuery) ||
                           this.normalize(unit.name).includes(normalizedQuery);
      const matchesType = !type || unit.type === type;
      return matchesQuery && matchesType;
    });

    // Sortuj: najpierw te które zaczynają się od query, potem reszta
    filtered.sort((a, b) => {
      const aStarts = this.normalize(a.name).startsWith(normalizedQuery);
      const bStarts = this.normalize(b.name).startsWith(normalizedQuery);
      
      if (aStarts && !bStarts) return -1;
      if (!aStarts && bStarts) return 1;
      return a.name.localeCompare(b.name, 'pl');
    });

    return filtered.slice(0, 20); // Zwróć maksymalnie 20 wyników
  }

  /**
   * Czeka na załadowanie danych
   */
  private waitForLoad(): Promise<void> {
    return new Promise((resolve) => {
      const checkInterval = setInterval(() => {
        if (!this.loading || this.loaded) {
          clearInterval(checkInterval);
          resolve();
        }
      }, 100);
      
      // Timeout po 10 sekundach
      setTimeout(() => {
        clearInterval(checkInterval);
        resolve();
      }, 10000);
    });
  }

  /**
   * Pobiera wszystkie jednostki określonego typu
   */
  getByType(type: TerritorialUnit['type']): TerritorialUnit[] {
    return this.units.filter(unit => unit.type === type);
  }

  /**
   * Generuje sugerowany email kontaktowy dla jednostki terytorialnej
   * Preferuje rzeczywisty email z bazy danych, jeśli dostępny
   */
  generateContactEmail(unit: TerritorialUnit): string {
    // Jeśli jednostka ma rzeczywisty email w bazie, użyj go
    if (unit.email && unit.email.trim()) {
      // Jeśli są multiple emaile (oddzielone średnikiem lub przecinkiem), użyj pierwszego
      const firstEmail = unit.email.split(/[;,]/)[0].trim();
      return firstEmail;
    }

    // Fallback: generuj email na podstawie nazwy
    let nazwa = unit.name
      .replace(/^Województwo\s+/i, '')
      .replace(/^Powiat\s+(m\.\s+)?/i, '')
      .replace(/^Gmina\s+/i, '')
      .replace(/^Miasto\s+/i, '');

    // Normalizuj do slug (bez polskich znaków, małe litery, myślniki zamiast spacji)
    const slug = this.normalize(nazwa)
      .replace(/\s+/g, '-')
      .replace(/[^a-z0-9-]/g, '');

    // Generuj email w zależności od typu
    switch (unit.type) {
      case 'wojewodztwo':
        return `kontakt@${slug}.uw.gov.pl`;
      case 'powiat':
        return `starostwo@${slug}.pl`;
      case 'miasto':
        return `urzad@um.${slug}.pl`;
      case 'gmina':
        return `ug@${slug}.pl`;
      default:
        return `kontakt@${slug}.pl`;
    }
  }

  /**
   * Normalizuje tekst dla porównania (małe litery, bez polskich znaków)
   */
  private normalize(text: string): string {
    return text
      .toLowerCase()
      .replace(/ą/g, 'a')
      .replace(/ć/g, 'c')
      .replace(/ę/g, 'e')
      .replace(/ł/g, 'l')
      .replace(/ń/g, 'n')
      .replace(/ó/g, 'o')
      .replace(/ś/g, 's')
      .replace(/ź/g, 'z')
      .replace(/ż/g, 'z');
  }
}
