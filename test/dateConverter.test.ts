import { expect } from 'chai';
import { ethers } from 'hardhat';
import { DateConverterMock } from '../typechain-types';

describe('DateConverter', function () {
    let dateConverter: DateConverterMock;

    before(async function () {
        const DateConverterFactory = await ethers.getContractFactory('DateConverterMock');
        dateConverter = await DateConverterFactory.deploy();
    });

    describe('YYMMDDhhmmssZ format (13 bytes)', function () {
        it('should convert valid date (2000s)', async function () {
            // 2023-06-15 14:30:45
            const timestamp = await dateConverter.convertDateToUnixTimestamp('230615143045Z');
            const expectedDate = new Date('2023-06-15T14:30:45Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));
        });

        it('should convert valid date (1900s)', async function () {
            // 1999-12-31 23:59:59
            const timestamp = await dateConverter.convertDateToUnixTimestamp('991231235959Z');
            const expectedDate = new Date('1999-12-31T23:59:59Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));
        });

        it('should handle leap year (29 Feb)', async function () {
            // 2020-02-24 00:00:00
            const timestamp0 = await dateConverter.convertDateToUnixTimestamp('200224000000Z');
            const expectedDate0 = new Date('2020-02-24T00:00:00Z');
            expect(timestamp0).to.equal(Math.floor(expectedDate0.getTime() / 1000));

            // 2020-02-29 00:00:00
            const timestamp = await dateConverter.convertDateToUnixTimestamp('200229000000Z');
            const expectedDate = new Date('2020-02-29T00:00:00Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));

            // 2020-03-05 00:00:00
            const timestamp2 = await dateConverter.convertDateToUnixTimestamp('200305000000Z');
            const expectedDate2 = new Date('2020-03-05T00:00:00Z');
            expect(timestamp2).to.equal(Math.floor(expectedDate2.getTime() / 1000));
        });

        it('should reject invalid day', async function () {
            // 2023-02-30 (invalid)
            await expect(dateConverter.convertDateToUnixTimestamp('230230000000Z')).to.be.rejectedWith('Invalid day');
        });

        it('should reject invalid month', async function () {
            // 2023-13-01 (invalid)
            await expect(dateConverter.convertDateToUnixTimestamp('231301000000Z')).to.be.rejectedWith('Invalid month');
        });

        it('should reject invalid hour', async function () {
            // 2023-01-01 24:00:00 (invalid)
            await expect(dateConverter.convertDateToUnixTimestamp('230101240000Z')).to.be.rejectedWith('Invalid hours');
        });

        it('should reject invalid minute', async function () {
            // 2023-01-01 00:60:00 (invalid)
            await expect(dateConverter.convertDateToUnixTimestamp('230101006000Z')).to.be.rejectedWith('Invalid minutes');
        });

        it('should reject invalid second', async function () {
            // 2023-01-01 00:00:60 (invalid)
            await expect(dateConverter.convertDateToUnixTimestamp('230101000060Z')).to.be.rejectedWith('Invalid seconds');
        });

        it('should reject non-numeric characters', async function () {
            // Contains 'A' instead of digit
            await expect(dateConverter.convertDateToUnixTimestamp('23A615143045Z')).to.be.rejectedWith('ASCII symbol is not a number');
        });
    });

    describe('YYYYMMDDhhmmss.mmmZ format (19 bytes)', function () {
        it('should convert valid date (milliseconds ignored)', async function () {
            // 2023-06-15 14:30:45.123
            const timestamp = await dateConverter.convertDateToUnixTimestamp('20230615143045.123Z');
            const expectedDate = new Date('2023-06-15T14:30:45Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));
        });

        it('should handle year 1970 (Unix epoch)', async function () {
            // 1970-01-01 00:00:00.000
            const timestamp = await dateConverter.convertDateToUnixTimestamp('19700101000000.000Z');
            expect(timestamp).to.equal(0);
        });

        it('should handle year 2100 (non-leap century)', async function () {
            // 2100-01-20 12:00:00.000 (2100 is not a leap year)
            const timestamp0 = await dateConverter.convertDateToUnixTimestamp('21000120120000.000Z');
            const expectedDate0 = new Date('2100-01-20T12:00:00Z');
            expect(timestamp0).to.equal(Math.floor(expectedDate0.getTime() / 1000));

            // 2100-02-29 12:00:00.000 (2100 is not a leap year)
            await expect(dateConverter.convertDateToUnixTimestamp('21000229120000.000Z')).to.be.rejectedWith('Invalid day');

            // 2100-03-01 12:00:00.000 (2100 is not a leap year)
            const timestamp2 = await dateConverter.convertDateToUnixTimestamp('21000301120000.000Z');
            const expectedDate2 = new Date('2100-03-01T12:00:00Z');
            expect(timestamp2).to.equal(Math.floor(expectedDate2.getTime() / 1000));
        });

        it('should handle year 2000 (leap century)', async function () {
            // 2000-02-25 12:00:00.000 (2000 is a leap year)
            const timestamp0 = await dateConverter.convertDateToUnixTimestamp('20000225120000.000Z');
            const expectedDate0 = new Date('2000-02-25T12:00:00Z');
            expect(timestamp0).to.equal(Math.floor(expectedDate0.getTime() / 1000));

            // 2000-02-29 12:00:00.000 (2000 is a leap year)
            const timestamp = await dateConverter.convertDateToUnixTimestamp('20000229120000.000Z');
            const expectedDate = new Date('2000-02-29T12:00:00Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));

            // 2000-03-05 12:00:00.000 (2000 is a leap year)
            const timestamp2 = await dateConverter.convertDateToUnixTimestamp('20000305120000.000Z');
            const expectedDate2 = new Date('2000-03-05T12:00:00Z');
            expect(timestamp2).to.equal(Math.floor(expectedDate2.getTime() / 1000));
        });

        it('should reject invalid format (missing dot)', async function () {
            // Missing dot before milliseconds
            await expect(dateConverter.convertDateToUnixTimestamp('20230615143045123Z')).to.be.rejectedWith(
                'Only 13 and 19 bytes are supported date formats'
            );
        });

        it('should reject year before 1970', async function () {
            // 1969-12-31 23:59:59.999
            await expect(dateConverter.convertDateToUnixTimestamp('19691231235959.999Z')).to.be.rejectedWith('Year must be >= 1970');
        });
    });

    describe('Edge cases', function () {
        it('should reject unsupported format length', async function () {
            // 14 bytes (invalid length)
            await expect(dateConverter.convertDateToUnixTimestamp('20230615143045Z')).to.be.rejectedWith(
                'Only 13 and 19 bytes are supported date formats'
            );
        });

        it('should handle maximum values (23:59:59)', async function () {
            // 1999-12-31 23:59:59
            const timestamp = await dateConverter.convertDateToUnixTimestamp('991231235959Z');
            const expectedDate = new Date('1999-12-31T23:59:59Z');
            expect(timestamp).to.equal(Math.floor(expectedDate.getTime() / 1000));
        });

        it('should handle minimum values (00:00:00)', async function () {
            // 1970-01-01 00:00:00
            const timestamp = await dateConverter.convertDateToUnixTimestamp('700101000000Z');
            expect(timestamp).to.equal(0);
        });
    });

    describe('Utility functions', function () {
        it('should correctly identify leap years', async function () {
            expect(await dateConverter.isLeapYear(2000)).to.equal(true); // Leap century
            expect(await dateConverter.isLeapYear(2004)).to.equal(true); // Leap year
            expect(await dateConverter.isLeapYear(2100)).to.equal(false); // Non-leap century
            expect(await dateConverter.isLeapYear(2001)).to.equal(false); // Non-leap year
        });

        it('should return correct days in month', async function () {
            // February in leap year
            expect(await dateConverter.getDaysInMonth(2, 2020)).to.equal(29);
            // February in non-leap year
            expect(await dateConverter.getDaysInMonth(2, 2021)).to.equal(28);
            // 30-day months
            expect(await dateConverter.getDaysInMonth(4, 2023)).to.equal(30);
            expect(await dateConverter.getDaysInMonth(6, 2023)).to.equal(30);
            expect(await dateConverter.getDaysInMonth(9, 2023)).to.equal(30);
            expect(await dateConverter.getDaysInMonth(11, 2023)).to.equal(30);
            // 31-day months
            expect(await dateConverter.getDaysInMonth(1, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(3, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(5, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(7, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(8, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(10, 2023)).to.equal(31);
            expect(await dateConverter.getDaysInMonth(12, 2023)).to.equal(31);
        });
    });
});
